package comms

import (
	"encoding/json"
)

func Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func Decode[T any](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	return result, err
}
