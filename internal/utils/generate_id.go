package utils

import (
	"crypto/rand"
	"fmt"
)

func GenerateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%X-%X", b[0:4], b[4:8])
}
