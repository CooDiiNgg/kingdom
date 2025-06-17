package main

import (
	"flag"
	"log"
	"os"

	agent "kingdom/internal/agents"
)

func main() {
	var (
		url      = flag.String("url", getenv("C2_URL", "http://127.0.0.1:8000"), "C2 base URL (scheme://host:port)")
		clientID = flag.String("client", getenv("CLIENT_ID", ""), "Client ID")
		agentID  = flag.String("agent", getenv("AGENT_ID", ""), "Agent ID")
	)
	flag.Parse()

	if *clientID == "" || *agentID == "" {
		log.Fatal("client and agent IDs must be provided via flags or env vars")
	}

	a := agent.New(*url, *clientID, *agentID)
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
