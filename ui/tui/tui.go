package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	cl "kingdom/internal/clients"
	commstypes "kingdom/internal/comms/comms_types"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type uiMode int

const (
	modeList uiMode = iota
	modeTask
	modeAgent
)

type agentItem struct{ ref cl.AgentRef }

func (a agentItem) Title() string       { return fmt.Sprintf("%s/%s", a.ref.ClientID, a.ref.AgentID) }
func (a agentItem) Description() string { return "" }
func (a agentItem) FilterValue() string { return a.Title() }

type errMsg struct{ error }

type agentsMsg []list.Item

type model struct {
	sdk *cl.Client

	list list.Model

	mode        uiMode
	prompt      textinput.Model
	promptLabel string

	err  error
	info string
}

func newModel(sdk *cl.Client) model {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = fmt.Sprintf("Agents for %s  (r:refresh  t:task  a:new agent  q:quit)", sdk.ID)

	ti := textinput.New()
	ti.CharLimit = 128

	return model{
		sdk:    sdk,
		list:   l,
		mode:   modeList,
		prompt: ti,
	}
}

func (m model) fetchAgents() tea.Cmd {
	return func() tea.Msg {
		refs, err := m.sdk.ListAgents()
		if err != nil {
			return errMsg{err}
		}
		items := make([]list.Item, 0, len(refs))
		for _, r := range refs {
			if r.ClientID == m.sdk.ID {
				items = append(items, agentItem{ref: r})
			}
		}
		return agentsMsg(items)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.fetchAgents(), tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
		return m, nil

	case agentsMsg:
		m.list.SetItems(msg)
		if len(msg) > 0 {
			m.list.Select(0)
		}
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case modeList:
			return m.updateList(msg)
		case modeTask, modeAgent:
			return m.updatePrompt(msg)
		}
	}
	return m, nil
}

func (m model) updateList(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "r":
		m.info = ""
		return m, m.fetchAgents()
	case "t":
		if len(m.list.Items()) == 0 {
			return m, nil
		}
		m.mode = modeTask
		m.promptLabel = "Command [args…]: "
		m.prompt.Reset()
		m.prompt.Focus()
		return m, nil
	case "a":
		m.mode = modeAgent
		m.promptLabel = "Platform (e.g. windows/amd64): "
		m.prompt.Reset()
		m.prompt.Focus()
		return m, nil
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(k)
	return m, cmd
}

func (m model) updatePrompt(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "enter":
		input := strings.TrimSpace(m.prompt.Value())
		if input == "" {
			m.mode = modeList
			m.prompt.Blur()
			return m, nil
		}

		switch m.mode {
		case modeTask:
			if sel, ok := m.list.SelectedItem().(agentItem); ok {
				cmd, args := splitOnce(input)
				task := &commstypes.Task{ID: fmt.Sprintf("%d", time.Now().UnixNano()), Command: cmd, Args: args}
				if err := m.sdk.QueueTask(sel.ref.ClientID, sel.ref.AgentID, task); err != nil {
					m.err = err
					m.info = ""
				} else {
					m.info = fmt.Sprintf("Queued %s %s", cmd, args)
					m.err = nil
				}
			}
			m.mode = modeList
			m.prompt.Blur()
			return m, m.fetchAgents()

		case modeAgent:
			platform := input
			if platform == "" {
				platform = "windows/amd64"
			}
			resp, err := m.sdk.CreateAgent(platform)
			if err != nil {
				m.err = err
				m.info = ""
			} else {
				_ = os.WriteFile(resp.FileName, []byte(resp.FileContent), 0o755)
				m.err = nil
				m.info = fmt.Sprintf("New agent %s created. Bootstrap script saved as %s.\nRun:\n%s", resp.AgentID, resp.FileName, strings.TrimSpace(resp.FileContent))
			}
			m.mode = modeList
			m.prompt.Blur()
			return m, m.fetchAgents()
		}

	case "esc":
		m.mode = modeList
		m.prompt.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	m.prompt, cmd = m.prompt.Update(k)
	return m, cmd
}

func (m model) View() string {
	var b strings.Builder
	switch m.mode {
	case modeList:
		b.WriteString(m.list.View())
	default:
		b.WriteString(m.promptLabel)
		b.WriteString(m.prompt.View())
		b.WriteString("\n[enter] confirm  [esc] cancel\n")
	}
	if m.err != nil {
		b.WriteString("\n! ")
		b.WriteString(m.err.Error())
	} else if m.info != "" {
		b.WriteString("\n> ")
		b.WriteString(m.info)
	}
	return b.String()
}

func splitOnce(s string) (string, string) {
	if idx := strings.IndexRune(s, ' '); idx >= 0 {
		return s[:idx], strings.TrimSpace(s[idx+1:])
	}
	return s, ""
}

func main() {
	base := os.Getenv("C2_URL")
	if base == "" {
		base = "http://127.0.0.1:8000"
	}
	sdk := cl.New(base)
	if err := sdk.Register(); err != nil {
		log.Fatalf("register client: %v", err)
	}

	p := tea.NewProgram(newModel(sdk), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
