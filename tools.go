package vivo

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strconv"

	"github.com/google/uuid"
)

func int64toString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}

func hMACSHA256HEX(data string, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	hash := h.Sum(nil)
	return hash
}

func base64encode(s []byte) string {
	return base64.StdEncoding.EncodeToString(s)
}

func GenerateRequestID() string {
	s := uuid.New().String()
	return s
}

func GenerateSessionID() string {
	s := uuid.New().String()
	return s
}
