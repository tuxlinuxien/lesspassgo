package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"os"
	"strconv"

	"crypto/hmac"

	"golang.org/x/crypto/pbkdf2"
)

const (
	iterations = 8192
	keylen     = 32
)

func encryptLogin(login, password string) []byte {
	var out = pbkdf2.Key([]byte(password), []byte(login), iterations, keylen, sha256.New)
	return []byte(fmt.Sprintf("%x", out))
}

func renderPassword(encLogin []byte, site string, len, counter int, template string) string {
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

func main() {
	login := flag.String("login", "", "login")
	masterPassword := flag.String("password", "", "password")
	site := flag.String("site", "", "domain.com")
	length := flag.Int("length", 12, "generated password length")
	counter := flag.Int("counter", 1, "counter")
	lower := flag.Bool("l", false, "password contains [a-z]")
	upper := flag.Bool("u", false, "password contains [A-Z]")
	numeric := flag.Bool("n", false, "password contains [0-9]")
	special := flag.Bool("s", false, "password contains @&%?,=[]_:-+*$#!'^~;()/.")
	flag.Parse()

	var template = ""
	if *lower == true {
		template += "vc"
	}
	if *upper == true {
		template += "VC"
	}
	if *numeric == true {
		template += "n"
	}
	if *special == true {
		template += "s"
	}
	if template == "" {
		fmt.Println("You need to define a password format")
		os.Exit(-1)
	}
	encLogin := encryptLogin(*login, *masterPassword)
	fmt.Println(renderPassword(encLogin, *site, *length, *counter, template))
}
