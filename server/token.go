package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

var secret []byte

func init() {
	secret = uuid.NewV4().Bytes()
}

func signUser(id string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(id))
	return fmt.Sprintf("%s.%x", id, mac.Sum(nil))
}

func checkToken(token string) (*UserModel, error) {
	tokenInfo := strings.Split(token, ".")
	if len(tokenInfo) != 2 {
		return nil, errors.New("invalid token")
	}
	t := signUser(tokenInfo[0])
	if token != t {
		return nil, errors.New("invalid token")
	}
	return GetUserByID(tokenInfo[0])
}
