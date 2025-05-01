package comms

import (
	"encoding/json"
)

func encode_internal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func decode_internal[T any](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	return result, err
}

func Encode(v any) ([]byte, []byte, error) {
	key, err := generateKey()
	if err != nil {
		return nil, nil, err
	}
	encodedData, err := encode_internal(v)
	if err != nil {
		return nil, nil, err
	}
	encryptedData, err := encrypt(encodedData, key)
	if err != nil {
		return nil, nil, err
	}
	return encryptedData, key, nil
}

func Decode[T any](data []byte, key []byte) (T, error) {
	decodedData, err := decrypt(data, key)
	if err != nil {
		return *new(T), err
	}
	result, err := decode_internal[T](decodedData)
	if err != nil {
		return *new(T), err
	}
	return result, nil
}
