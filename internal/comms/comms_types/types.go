package commstypes

type Request struct {
	AgentID  string `json:"agent_id"`
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	IPAddr   string `json:"ipaddr"`
	Port     int    `json:"port"`
}

type Task struct {
	ID      string `json:"id"`
	Command string `json:"command"`
	Args    string `json:"args"`
}

type TaskResult struct {
	AgentID string `json:"agent_id"`
	TaskID  string `json:"task_id"`
	Status  string `json:"status"`
	Output  string `json:"output"`
	Error   string `json:"error"`
}
