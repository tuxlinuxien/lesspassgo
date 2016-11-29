// Package v1 provides core functions to build LessPass password.
package v1

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"strconv"

	"golang.org/x/crypto/pbkdf2"
)

const (
	iterations = 8192
	keylen     = 32
)

// EncryptLogin encrypts login with pbkdf2.
func EncryptLogin(login, password string) []byte {
	var out = pbkdf2.Key([]byte(password), []byte(login), iterations, keylen, sha256.New)
	return []byte(fmt.Sprintf("%x", out))
}

// RenderPassword returns the generated password.
func RenderPassword(encLogin []byte, site string, len, counter int, template string) string {
	derivedEncryptedLogin := deriveEncryptedLogin(encLogin, site, len, counter)
	return prettyPrint(derivedEncryptedLogin, template)
}

func createHmac(encLogin []byte, salt string) []byte {
	mac := hmac.New(sha256.New, encLogin)
	mac.Write([]byte(salt))
	return []byte(fmt.Sprintf("%x", mac.Sum(nil)))
}

func deriveEncryptedLogin(encLogin []byte, site string, length, counter int) []byte {
	var salt = site + strconv.Itoa(counter)
	return createHmac(encLogin, salt)[0:length]
}

func getPasswordChar(charType byte, index int) byte {
	var passwordsChars = map[byte]string{
		'V': "AEIOUY",
		'C': "BCDFGHJKLMNPQRSTVWXZ",
		'v': "aeiouy",
		'c': "bcdfghjklmnpqrstvwxz",
		'A': "AEIOUYBCDFGHJKLMNPQRSTVWXZ",
		'a': "AEIOUYaeiouyBCDFGHJKLMNPQRSTVWXZbcdfghjklmnpqrstvwxz",
		'n': "0123456789",
		's': "@&%?,=[]_:-+*$#!'^~;()/.",
		'x': "AEIOUYaeiouyBCDFGHJKLMNPQRSTVWXZbcdfghjklmnpqrstvwxz0123456789@&%?,=[]_:-+*$#!'^~;()/.",
	}

	var passwordChar = passwordsChars[charType]
	return passwordChar[index%len(passwordChar)]
}

func getCharType(template string, index int) byte {
	return template[index%len(template)]
}

func prettyPrint(hash []byte, template string) string {
	var out = ""
	for i, c := range hash {
		tplStr := getCharType(template, i)
		out += string(getPasswordChar(tplStr, int(c)))
	}
	return out
}
