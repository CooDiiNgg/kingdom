package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	scheduler "kingdom/internal/c2"
	commstypes "kingdom/internal/comms/comms_types"
	"kingdom/internal/listeners"
	listenerstypes "kingdom/internal/listeners/listeners_types"
	"kingdom/internal/storage"

	"github.com/google/uuid"
)

type agentRef struct {
	ClientID string `json:"client_id"`
	AgentID  string `json:"agent_id"`
}

type agentMeta struct {
	ref      agentRef
	platform string
	port     int
}

var (
	clients = struct {
		sync.RWMutex
		m map[string]struct{}
	}{m: make(map[string]struct{})}

	agents = struct {
		sync.RWMutex
		m map[string]map[string]*agentMeta
	}{m: make(map[string]map[string]*agentMeta)}
)

func main() {
	apiAddr := flag.String("addr", getenv("API_ADDR", ":8000"), "[host]:port for the client-facing API server")
	dbPath := flag.String("db", getenv("DB_PATH", "kingdom.db"), "SQLite storage path for session keys")
	// pruneIntvl := flag.Duration("prune", 30*time.Minute, "interval to drop inactive agents from scheduler")
	flag.Parse()

	store, err := storage.NewSQLite(*dbPath)
	if err != nil {
		log.Fatalf("storage: %v", err)
	}
	storage.SetDefaultProvider(store)
	defer store.Close()

	go scheduler.ListAgents()
	// go func() { scheduler.Default().Prune(*pruneIntvl) }()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/clients/", handleClients)
	mux.HandleFunc("/api/clients", handleAgents)
	mux.HandleFunc("/api/agents", handleAgents)
	mux.HandleFunc("/api/agents/", handleAgentTasks)

	srv := &http.Server{
		Addr:    *apiAddr,
		Handler: mux,
	}

	go func() {
		log.Printf("[api] ready on http://%s", *apiAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("api: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func handleClients(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	parts := pathParts(r.URL.Path)
	if len(parts) != 3 || parts[0] != "api" || parts[1] != "clients" {
		if len(parts) == 4 && parts[0] == "api" && parts[1] == "clients" && parts[3] == "agents" {
			handleAgents(w, r)
			return
		}
		http.NotFound(w, r)
		return
	}
	clientID := parts[2]

	clients.Lock()
	defer clients.Unlock()
	if _, exists := clients.m[clientID]; exists {
		http.Error(w, "client already registered", http.StatusConflict)
		return
	}
	clients.m[clientID] = struct{}{}
	w.WriteHeader(http.StatusCreated)
}

// GET  /api/agents                       – list all agents
// POST /api/clients/{client}/agents      – create a new agent under client
func handleAgents(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/agents" {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		listAgents(w)
		return
	}

	// Expect: /api/clients/{clientID}/agents
	parts := pathParts(r.URL.Path)
	if len(parts) == 4 && parts[0] == "api" && parts[1] == "clients" && parts[3] == "agents" {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		createAgent(w, r, parts[2])
		return
	}

	http.NotFound(w, r)
}

// POST /api/agents/{client}/{agent}/tasks – enqueue task for an agent
func handleAgentTasks(w http.ResponseWriter, r *http.Request) {
	parts := pathParts(r.URL.Path)
	if len(parts) != 5 || parts[0] != "api" || parts[1] != "agents" || parts[4] != "tasks" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	clientID, agentID := parts[2], parts[3]
	var t commstypes.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := scheduler.QueueTask(clientID, agentID, &t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func listAgents(w http.ResponseWriter) {
	out := make([]agentRef, 0)
	agents.RLock()
	for _, perClient := range agents.m {
		for _, meta := range perClient {
			out = append(out, meta.ref)
		}
	}
	agents.RUnlock()
	respondJSON(w, out, http.StatusOK)
}

func createAgent(w http.ResponseWriter, r *http.Request, clientID string) {
	if !clientRegistered(clientID) {
		http.Error(w, "unknown client", http.StatusNotFound)
		return
	}

	var in struct {
		Platform string `json:"platform"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	agentID := uuid.NewString()

	port, err := freePort()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addr := listenerstypes.ListenerAddr{Address: "", Port: port}
	lst, err := listeners.NewHttpListener(addr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := listeners.RegisterListener(clientID, agentID, lst); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	agents.Lock()
	if _, ok := agents.m[clientID]; !ok {
		agents.m[clientID] = make(map[string]*agentMeta)
	}
	agents.m[clientID][agentID] = &agentMeta{ref: agentRef{ClientID: clientID, AgentID: agentID}, platform: in.Platform, port: port}
	agents.Unlock()

	baseURL := fmt.Sprintf("http://%s:%d", hostOnly(r.Host), port)
	fileContent := fmt.Sprintf("#!/usr/bin/env sh\nCLIENT_ID=%s AGENT_ID=%s C2_URL=%s ./agent\n", clientID, agentID, baseURL)

	respondJSON(w, map[string]any{
		"agent_id":     agentID,
		"file_content": fileContent,
		"file_name":    fmt.Sprintf("agent_%s_bootstrap.sh", agentID),
	}, http.StatusCreated)
}

func respondJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func clientRegistered(id string) bool {
	clients.RLock()
	_, ok := clients.m[id]
	clients.RUnlock()
	return ok
}

func pathParts(p string) []string {
	segs := strings.Split(strings.TrimPrefix(p, "/"), "/")
	out := segs[:0]
	for _, s := range segs {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func freePort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func hostOnly(h string) string {
	if idx := strings.IndexRune(h, ':'); idx >= 0 {
		return h[:idx]
	}
	if h == "" {
		return "127.0.0.1"
	}
	return h
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
