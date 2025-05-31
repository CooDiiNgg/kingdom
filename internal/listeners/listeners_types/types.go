package listenerstypes

type Listener interface {
	Start(clientID string) error
	Stop() error
}