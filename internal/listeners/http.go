package listeners

import (
	"net/http"
)

type HTTPListener struct {
	address string
	port    int
	clientID string
	server *http.Server
}

func (h *HTTPListener) Start(clientID string) error {
	h.clientID = clientID
	mux := http.NewServeMux()

	mux.HandleFunc(fmt.Sprintf("/client/%s", h.clientID), func(w http.ResponseWriter, r *http.Request) {
		resp, err := HandleRequest(h.clientID, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	if h.address == "" {
		h.address = "localhost"
	}
	if h.port == 0 {
		h.port = 8080
	}

	h.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", h.address, h.port),
		Handler: mux,
	}

	go h.server.ListenAndServe()
	return nil
}

func (h *HTTPListener) Stop() error {
	if h.server != nil {
		return h.server.Close()
	}
	return nil
}

func (h *HTTPListener) Configure(addr ListenerAddr) error {
	if addr.Address != "" {
		h.address = addr.Address
	}
	if addr.Port != 0 {
		h.port = addr.Port
	}
	return nil
}