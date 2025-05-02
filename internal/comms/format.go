package comms

import (
	"encoding/json"
	"errors"
)

func encode_internal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func decode_internal[T any](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	return result, err
}

func Encode(v any, key_arg ...[]byte) ([]byte, []byte, error) {
	var key []byte
	var err error
	if len(key_arg) == 1 {
		if len(key_arg[0]) != 32 {
			return nil, nil, errors.New("Key must be 32 bytes")
		}
		key = key_arg[0]
		err = nil
	} else if len(key_arg) == 0 {
		key, err = generateKey()
	} else {
		err = errors.New("Too many arguments")
	}

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
	if len(key) != 32 {
		return *new(T), errors.New("Key must be 32 bytes")
	}
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
