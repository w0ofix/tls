package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
