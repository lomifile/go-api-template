package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// RandomString Provides random string
func RandomString(n int) string {
	buf := make([]byte, n)
	_, err := rand.Read(buf)
	if err != nil {
		panic("rand.Read failed: " + err.Error())
	}
	return base64.RawURLEncoding.EncodeToString(buf)[:n]
}
