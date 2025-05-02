package unit

import (
	comms "kingdom/internal/comms/comms_types"
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

var CommsTestCases_DecodeInternal = []CommsTestCase_DecodeInternal{
	{
		Name:  "Happy path 1",
		Input: []byte(`{"agent_id":"agent1","hostname":"host1","os":"linux","ipaddr":"127.0.0.1","port":8080}`),
		Expected: &comms.Request{
			AgentID:  "agent1",
			Hostname: "host1",
			OS:       "linux",
			IPAddr:   "127.0.0.1",
			Port:     8080,
		},
		Error: false,
	},
	{
		Name:  "Happy path 2",
		Input: []byte(`{"id":"task1","command":"ls","args":"-l"}`),
		Expected: &comms.Task{
			ID:      "task1",
			Command: "ls",
			Args:    "-l",
		},
		Error: false,
	},
	{
		Name:  "Happy path 3",
		Input: []byte(`{"agent_id":"agent1","task_id":"task1","status":"success","output":"output","error":""}`),
		Expected: &comms.TaskResult{
			AgentID: "agent1",
			TaskID:  "task1",
			Status:  "success",
			Output:  "output",
			Error:   "",
		},
		Error: false,
	},
	{
		Name:     "Negative path 1",
		Input:    []byte(`{"agent_id"}`),
		Expected: nil,
		Error:    true,
	},
	{
		Name:     "Negative path 2",
		Input:    []byte("123"),
		Expected: nil,
		Error:    true,
	},
}

var CommsTestCases_Encrypt = []CommsTestCase_Encrypt{
	{
		Name:     "Happy path 1",
		Input:    []byte("Hello, World!"),
		Key:      []byte("12345678901234567890123456789012"),
		Expected: []byte{0x8a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f, 0x6a, 0x7b, 0x8c, 0x9d, 0xae, 0xbf, 0xd0, 0xe1, 0xf2, 0x03}, // Will change later to real values
		Error:    false,
	},
	{
		Name:     "Happy path 2",
		Input:    []byte(`{"agent_id":"agent1","hostname":"host1","os":"linux"}`),
		Key:      []byte("12345678901234567890123456789012"),
		Expected: []byte{0x8a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f, 0x6a, 0x7b, 0x8c, 0x9d, 0xae, 0xbf, 0xd0, 0xe1, 0xf2, 0x03}, // Will change later to real values
		Error:    false,
	},
	{
		Name:     "Negative path 1",
		Input:    []byte("Hello, World!"),
		Key:      []byte("1234567890123456789012345678901"), // 31 bytes
		Expected: nil,
		Error:    true,
	},
}

var CommsTestCases_Decrypt = []CommsTestCase_Decrypt{
	{
		Name:     "Happy path 1",
		Input:    []byte{0x8a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f, 0x6a, 0x7b, 0x8c, 0x9d, 0xae, 0xbf, 0xd0, 0xe1, 0xf2, 0x03}, // Will change later to real values
		Key:      []byte("12345678901234567890123456789012"),
		Expected: []byte("Hello, World!"),
		Error:    false,
	},
	{
		Name:     "Happy path 2",
		Input:    []byte{0x8a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f, 0x6a, 0x7b, 0x8c, 0x9d, 0xae, 0xbf, 0xd0, 0xe1, 0xf2, 0x03}, // Will change later to real values
		Key:      []byte("12345678901234567890123456789012"),
		Expected: []byte(`{"agent_id":"agent1","hostname":"host1","os":"linux"}`),
		Error:    false,
	},
	{
		Name:     "Negative path 1",
		Input:    []byte{0x8a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f, 0x6a, 0x7b, 0x8c, 0x9d, 0xae, 0xbf, 0xd0, 0xe1, 0xf2, 0x03}, // Will change later to real values
		Key:      []byte("1234567890123456789012345678901"),                                                              // 31 bytes
		Expected: nil,
		Error:    true,
	},
	{
		Name:     "Negative path 2",
		Input:    []byte("Hello, World!"),
		Key:      []byte("12345678901234567890123456789012"),
		Expected: nil,
		Error:    true,
	},
}

var CommsTestCases_Encode = []CommsTestCase_Encode{
	{
		Name: "Happy path 1",
		Input: &comms.Request{
			AgentID:  "agent1",
			Hostname: "host1",
			OS:       "linux",
			IPAddr:   "127.0.0.1",
			Port:     8080,
		},
		Expected: []byte{0x8a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f, 0x6a, 0x7b, 0x8c, 0x9d, 0xae, 0xbf, 0xd0, 0xe1, 0xf2, 0x03}, // Will change later to real values
		Key:      []byte("12345678901234567890123456789012"),
		Error:    false,
	},
	{
		Name: "Happy path 2",
		Input: &comms.Task{
			ID:      "task1",
			Command: "ls",
			Args:    "-l",
		},
		Expected: []byte{0x8a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f, 0x6a, 0x7b, 0x8c, 0x9d, 0xae, 0xbf, 0xd0, 0xe1, 0xf2, 0x03}, // Will change later to real values
		Key:      []byte("12345678901234567890123456789012"),
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
		Expected: []byte{0x8a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f, 0x6a, 0x7b, 0x8c, 0x9d, 0xae, 0xbf, 0xd0, 0xe1, 0xf2, 0x03}, // Will change later to real values
		Key:      []byte("12345678901234567890123456789012"),
		Error:    false,
	},
	{
		Name:     "Negative path 1",
		Input:    make(chan int),
		Expected: nil,
		Key:      []byte("12345678901234567890123456789012"),
		Error:    true,
	},
	{
		Name: "Negative path 2",
		Input: &comms.Task{
			ID:      "task2",
			Command: "ls",
			Args:    "-l",
		},
		Expected: nil,
		Key:      []byte("1234567890123456789012345678901"), // 31 bytes
		Error:    true,
	},
}

var CommsTestCases_Decode = []CommsTestCase_Decode{
	{
		Name:  "Happy path 1",
		Input: []byte{0x8a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f, 0x6a, 0x7b, 0x8c, 0x9d, 0xae, 0xbf, 0xd0, 0xe1, 0xf2, 0x03}, // Will change later to real values
		Key:   []byte("12345678901234567890123456789012"),
		Expected: &comms.Request{
			AgentID:  "agent1",
			Hostname: "host1",
			OS:       "linux",
			IPAddr:   "127.0.0.1",
			Port:     8080,
		},
		Error: false,
	},
	{
		Name:  "Happy path 2",
		Input: []byte{0x8a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f, 0x6a, 0x7b, 0x8c, 0x9d, 0xae, 0xbf, 0xd0, 0xe1, 0xf2, 0x03}, // Will change later to real values
		Key:   []byte("12345678901234567890123456789012"),
		Expected: &comms.Task{
			ID:      "task1",
			Command: "ls",
			Args:    "-l",
		},
		Error: false,
	},
	{
		Name:  "Happy path 3",
		Input: []byte{0x8a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f, 0x6a, 0x7b, 0x8c, 0x9d, 0xae, 0xbf, 0xd0, 0xe1, 0xf2, 0x03}, // Will change later to real values
		Key:   []byte("12345678901234567890123456789012"),
		Expected: &comms.TaskResult{
			AgentID: "agent1",
			TaskID:  "task1",
			Status:  "success",
			Output:  "output",
			Error:   "",
		},
		Error: false,
	},
	{
		Name:     "Negative path 1",
		Input:    []byte{0x8a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f, 0x6a, 0x7b, 0x8c, 0x9d, 0xae, 0xbf, 0xd0, 0xe1, 0xf2, 0x03}, // Will change later to real values
		Key:      []byte("1234567890123456789012345678901"),                                                              // 31 bytes
		Expected: nil,
		Error:    true,
	},
	{
		Name:     "Negative path 2",
		Input:    []byte("Hello, World!"),
		Key:      []byte("12345678901234567890123456789012"),
		Expected: nil,
		Error:    true,
	},
	{
		Name:     "Negative path 3",
		Input:    []byte{0x8a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f, 0x6a, 0x7b, 0x8c, 0x9d, 0xae, 0xbf, 0xd0, 0xe1, 0xf2, 0x03}, // Will change later to real values - they should decrypt to a string
		Key:      []byte("12345678901234567890123456789012"),
		Expected: nil,
		Error:    true,
	},
}
