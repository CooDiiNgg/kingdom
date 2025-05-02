package unit

import (
	"kingdom/internal/comms"
)

type CommsTestCase_EncodeInternal struct {
	Name     string
	Input    any
	Expected []byte
	Error    bool
}

type CommsTestCase_DecodeInternal struct {
	Name     string
	Input    []byte
	Expected any
	Error    bool
}

type CommsTestCase_Encrypt struct {
	Name     string
	Input    []byte
	Key      []byte
	Expected []byte
	Error    bool
}

type CommsTestCase_Decrypt struct {
	Name     string
	Input    []byte
	Key      []byte
	Expected []byte
	Error    bool
}

type CommsTestCase_Encode struct {
	Name     string
	Input    any
	Expected []byte
	Key      []byte
	Error    bool
}

type CommsTestCase_Decode struct {
	Name     string
	Input    []byte
	Key      []byte
	Expected any
	Error    bool
}

var CommsTestCases_EncodeInternal = []CommsTestCase_EncodeInternal{
	{
		Name: "Happy path 1",
		Input: &comms.Request{
			AgentID:  "agent1",
			Hostname: "host1",
			OS:       "linux",
			IPAddr:   "127.0.0.1",
			Port:     8080,
		},
		Expected: []byte(`{"agent_id":"agent1","hostname":"host1","os":"linux","ipaddr":"127.0.1","port":8080}`),
		Error:    false,
	},
	{
		Name: "Happy path 2",
		Input: &comms.Task{
			ID:      "task1",
			Command: "ls",
			Args:    "-l",
		},
		Expected: []byte(`{"id":"task1","command":"ls","args":"-l"}`),
		Error:    false,
	},
	{
		Name: "Happy path 3",
		Input: &comms.TaskResult{
			AgentID: "agent1",
			TaskID:  "task1",
			Status:  "success",
			Output:  "output",
			Error:   "",
		},
		Expected: []byte(`{"agent_id":"agent1","task_id":"task1","status":"success","output":"output","error":""}`),
		Error:    false,
	},
	{
		Name:     "Negative path 1",
		Input:    make(chan int),
		Expected: nil,
		Error:    true,
	},
	{
		Name:     "Negative path 2",
		Input:    123,
		Expected: nil,
		Error:    true,
	},
	{
		Name:     "Negative path 3",
		Input:    []int{1, 2, 3},
		Expected: nil,
		Error:    true,
	},
}
