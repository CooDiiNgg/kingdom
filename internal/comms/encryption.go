package comms

import (
	"crypto/aes"
	"crypto/rand"
	"errors"
)

func encrypt(data []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("Key must be 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(data))
	block.Encrypt(ciphertext, data)

	return ciphertext, nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("Key must be 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	block.Decrypt(plaintext, ciphertext)

	return plaintext, nil
}

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
