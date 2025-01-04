package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func pkcs5Unpad(data []byte) ([]byte, error) {
	padding := int(data[len(data)-1])
	if padding > len(data) {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:len(data)-padding], nil
}

func decryptAES(iv, key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(ciphertext))

	mode.CryptBlocks(plaintext, ciphertext)
	plaintext, err = pkcs5Unpad(plaintext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
