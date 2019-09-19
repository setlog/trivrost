package misc

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MustGetRandomHexString(byteCount int) string {
	bytes := make([]byte, byteCount)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(fmt.Sprintf("Could not rand.Read(): %v", err))
	}
	return hex.EncodeToString(bytes)
}
