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

func Encode(v any, key_and_iv ...[]byte) ([]byte, []byte, []byte, error) {
	var key []byte
	var iv []byte
	var err error
	if len(key_and_iv) == 2 {
		if len(key_and_iv[0]) == 32 && len(key_and_iv[1]) == 16 {
			key = key_and_iv[0]
			iv = key_and_iv[1]
			err = nil
		} else if len(key_and_iv[0]) == 16 && len(key_and_iv[1]) == 32 {
			iv = key_and_iv[0]
			key = key_and_iv[1]
			err = nil
		} else {
			err = errors.New("Key must be 32 bytes and IV must be 16 bytes")
		}
	} else if len(key_and_iv) == 1 {
		if len(key_and_iv[0]) == 32 {
			key = key_and_iv[0]
			iv, err = generateIV()
			if err != nil {
				return nil, nil, nil, err
			}
		} else if len(key_and_iv[0]) == 16 {
			iv = key_and_iv[0]
			key, err = generateKey()
			if err != nil {
				return nil, nil, nil, err
			}
		} else {
			err = errors.New("Key must be 32 bytes or IV must be 16 bytes")
		}
	} else if len(key_and_iv) == 0 {
		key, err = generateKey()
		if err != nil {
			return nil, nil, nil, err
		}
		iv, err = generateIV()
		if err != nil {
			return nil, nil, nil, err
		}
	} else {
		err = errors.New("Too many arguments")
	}

	if err != nil {
		return nil, nil, nil, err
	}
	encodedData, err := encode_internal(v)
	if err != nil {
		return nil, nil, nil, err
	}
	encryptedData, err := encrypt(encodedData, key, iv)
	if err != nil {
		return nil, nil, nil, err
	}
	return encryptedData, key, iv, nil
}

func Decode[T any](data []byte, key []byte, iv []byte) (T, error) {
	if len(key) != 32 {
		return *new(T), errors.New("Key must be 32 bytes")
	}
	if len(iv) != 16 {
		return *new(T), errors.New("IV must be 16 bytes")
	}
	if len(data) == 0 {
		return *new(T), errors.New("data is empty")
	}
	decodedData, err := decrypt(data, key, iv)
	if err != nil {
		return *new(T), err
	}
	result, err := decode_internal[T](decodedData)
	if err != nil {
		return *new(T), err
	}
	return result, nil
}
