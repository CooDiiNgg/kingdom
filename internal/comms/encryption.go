package comms

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
)

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	padding := data[len(data)-1]
	if int(padding) > len(data) {
		return nil, errors.New("padding size is larger than data size")
	}
	return data[:len(data)-int(padding)], nil
}

func encrypt(data []byte, key []byte, iv []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("Key must be 32 bytes")
	}
	if len(iv) != aes.BlockSize {
		return nil, errors.New("IV must be 16 bytes")
	}
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	data = pkcs7Pad(data, block.BlockSize())
	if len(data)%block.BlockSize() != 0 {
		return nil, errors.New("data is not a multiple of block size")
	}
	ciphertext := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, data)
	if len(ciphertext) != len(data) {
		return nil, errors.New("ciphertext length does not match data length")
	}
	return ciphertext, nil
}

func decrypt(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("Key must be 32 bytes")
	}
	if len(iv) != aes.BlockSize {
		return nil, errors.New("IV must be 16 bytes")
	}
	if len(ciphertext) == 0 {
		return nil, errors.New("ciphertext is empty")
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of block size")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)
	if len(plaintext) != len(ciphertext) {
		return nil, errors.New("plaintext length does not match ciphertext length")
	}
	plaintext, err = pkcs7Unpad(plaintext)
	if err != nil {
		return nil, err
	}
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

func generateIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	_, err := rand.Read(iv)
	if err != nil {
		return nil, err
	}
	return iv, nil
}
