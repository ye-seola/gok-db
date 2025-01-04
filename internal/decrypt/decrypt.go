package decrypt

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

var (
	kdbPrefixes = []string{
		"", "", "12", "24", "18", "30", "36", "12", "48", "7", "35", "40", "17",
		"23", "29", "isabel", "kale", "sulli", "van", "merry", "kyle", "james",
		"maddux", "tony", "hayden", "paul", "elijah", "dorothy", "sally",
		"bran", "extr.ursra", "veil",
	}
	kdbPW = []byte{0, 22, 0, 8, 0, 9, 0, 111, 0, 2, 0, 23, 0, 43, 0, 8, 0, 33, 0, 33, 0, 10, 0, 16, 0, 3, 0, 3, 0, 7, 0, 6, 0, 0}
	kdbIV = []byte{15, 8, 1, 0, 25, 71, 37, 220, 21, 245, 23, 224, 225, 21, 12, 53}
)

func generateSalt(userId uint64, encType uint32) ([]byte, error) {
	if int(encType) >= len(kdbPrefixes) {
		return nil, fmt.Errorf("invalid EncType")
	}

	salt := make([]byte, 16)
	saltSource := kdbPrefixes[encType] + strconv.FormatUint(userId, 10)
	copy(salt, saltSource)

	return salt, nil
}

func deriveKey(userId uint64, encType uint32) ([]byte, error) {
	salt, err := generateSalt(userId, encType)
	if err != nil {
		return nil, err
	}

	return sha1Pbkdf(salt, kdbPW, 2, 32), nil
}

func Decrypt(userId uint64, data string, encType uint32) (string, error) {
	data = strings.TrimSpace(data)
	if data == "" {
		return "", nil
	}

	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	key, err := deriveKey(userId, encType)
	if err != nil {
		return "", err
	}

	plain, err := decryptAES(kdbIV, key, ciphertext)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
