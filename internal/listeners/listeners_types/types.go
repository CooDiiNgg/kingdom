package listenerstypes

type Listener interface {
	Start(clientID string, agentID string) error
	Stop() error
	Configure(addr ListenerAddr) error
}

type ListenerAddr struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}