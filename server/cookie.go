package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

var secret []byte

func init() {
	secret = uuid.NewV4().Bytes()
}

func signUser(email string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(email))
	return fmt.Sprintf("%s.%x", email, mac.Sum(nil))
}
