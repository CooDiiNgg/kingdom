package listeners

var RegisteredListeners = make(map[string]map[string]Listener)

func RegisterListener(clientID string, agentID string, listener Listener) {
	if _, exists := RegisteredListeners[clientID]; !exists {
		RegisteredListeners[clientID] = make(map[string]Listener)
	}

	if _, exists := RegisteredListeners[clientID][agentID]; exists {
		return fmt.Errorf("Listener for client ID %s and agent ID %s already exists", clientID, agentID)
	}

	err := listener.Start(clientID, agentID)

	if err != nil {
		return fmt.Errorf("Failed to start listener for client ID %s and agent ID %s: %v", clientID, agentID, err)
	}

	RegisteredListeners[clientID][agentID] = listener
	return nil
}

func UnregisterListener(clientID string, agentID string) error {
	if _, exists := RegisteredListeners[clientID]; !exists {
		return fmt.Errorf("No listeners registered for client ID %s", clientID)
	}

	if _, exists := RegisteredListeners[clientID][agentID]; !exists {
		return fmt.Errorf("No listener registered for agent ID %s under client ID %s", agentID, clientID)
	}

	listener := RegisteredListeners[clientID][agentID]
	err := listener.Stop()
	if err != nil {
		return fmt.Errorf("Failed to stop listener for client ID %s and agent ID %s: %v", clientID, agentID, err)
	}

	delete(RegisteredListeners[clientID], agentID)
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