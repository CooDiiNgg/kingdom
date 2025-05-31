package listeners

var RegisteredListeners = make(map[string]Listener)

func RegisterListener(clientID string, listener Listener) {
	if _, exists := RegisteredListeners[clientID]; exists {
		return fmt.Errorf("Listener for client ID %s already exists", clientID)
	}
	
	err := listener.Start(clientID)
	if err != nil {
		return fmt.Errorf("Failed to start listener for client ID %s: %v", clientID, err)
	}

	RegisteredListeners[clientID] = listener

	return nil
}

func UnregisterListener(clientID string) error {
	listener, exists := RegisteredListeners[clientID]
	if !exists {
		return fmt.Errorf("Listener for client ID %s does not exist", clientID)
	}

	err := listener.Stop()
	if err != nil {
		return fmt.Errorf("Failed to stop listener for client ID %s: %v", clientID, err)
	}

	delete(RegisteredListeners, clientID)

	return nil
}

func NewHttpListener(addr ListenerAddr) (Listener, error) {
	listener := &HTTPListener{}
	err := listener.Configure(addr)
	if err != nil {
		return nil, fmt.Errorf("Failed to configure HTTP listener: %v", err)
	}
	return listener, nil
}