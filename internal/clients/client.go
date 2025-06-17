package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	commstypes "kingdom/internal/comms/comms_types"

	"github.com/google/uuid"
)

type AgentRef struct {
	ClientID string `json:"client_id"`
	AgentID  string `json:"agent_id"`
}

type Client struct {
	BaseURL string
	http    *http.Client
	ID      string
}

type newAgentRequest struct {
	Platform string `json:"platform"`
}

type newAgentResponse struct {
	AgentID     string `json:"agent_id"`
	FileContent string `json:"file_content"`
	FileName    string `json:"file_name"`
}

func New(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) Register() error {
	const maxRetries = 5
	for i := 0; i < maxRetries; i++ {
		candidate := uuid.NewString()
		uri := fmt.Sprintf("%s/api/clients/%s", c.BaseURL, candidate)
		req, _ := http.NewRequest(http.MethodPost, uri, nil)
		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
		switch resp.StatusCode {
		case http.StatusCreated, http.StatusOK:
			c.ID = candidate
			return nil
		case http.StatusConflict:
			continue
		default:
			return fmt.Errorf("register client: http %d", resp.StatusCode)
		}
	}
	return fmt.Errorf("failed to register client after %d attempts", maxRetries)
}

func (c *Client) ListAgents() ([]AgentRef, error) {
	uri := fmt.Sprintf("%s/api/agents", c.BaseURL)
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list agents: http %d", resp.StatusCode)
	}

	var out []AgentRef
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) QueueTask(clientID, agentID string, t *commstypes.Task) error {
	if t == nil {
		return fmt.Errorf("task is nil")
	}
	if t.ID == "" {
		return fmt.Errorf("task.ID cannot be empty")
	}

	uri := fmt.Sprintf("%s/api/agents/%s/%s/tasks", c.BaseURL, clientID, agentID)
	payload, err := json.Marshal(t)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("queue task: http %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) CreateAgent(platform string) (*newAgentResponse, error) {
	if c.ID == "" {
		return nil, fmt.Errorf("client not registered")
	}
	uri := fmt.Sprintf("%s/api/clients/%s/agents", c.BaseURL, c.ID)
	body, _ := json.Marshal(newAgentRequest{Platform: platform})
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create agent: http %d", resp.StatusCode)
	}
	var out newAgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
