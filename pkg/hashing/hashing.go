package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func HashImageName(imageName string) string {
	name := strings.Split(imageName, ".")[0]
	hash := sha256.New()
	hash.Write([]byte(name + time.Now().Format(time.RFC3339)))
	return hex.EncodeToString(hash.Sum(nil))
}
