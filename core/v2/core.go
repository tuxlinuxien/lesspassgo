// Package v2 provides core functions to build LessPass password.
package v2

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"golang.org/x/crypto/pbkdf2"
)

const (
	iterations = 100000
	keylen     = 32
)

var (
	charRules = map[string]string{
		"lowercase": "abcdefghijklmnopqrstuvwxyz",
		"uppercase": "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"numbers":   "0123456789",
		"symbols":   "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~",
	}
)

// PasswordProfile .
type PasswordProfile struct {
	Keylen     int
	Iterations int
	Counter    int
	Length     int
	Digest     string
	Rules      []string
}

// NewPasswordProfile returns the default password configuration.
func NewPasswordProfile() *PasswordProfile {
	return &PasswordProfile{
		keylen,
		iterations,
		1,
		16,
		"sha256",
		[]string{},
	}
}

// GeneratePassword generates v2 password.
func GeneratePassword(site, login, masterPassword string, pp *PasswordProfile) string {
	entropy := calcEntropy(site, login, masterPassword, pp)
	return renderPassword(entropy, pp)
}

func calcEntropy(site, login, masterPassword string, pp *PasswordProfile) []byte {
	var salt = site + login + fmt.Sprintf("%x", pp.Counter)
	var out []byte
	if pp.Digest == "sha256" {
		out = pbkdf2.Key([]byte(masterPassword), []byte(salt), pp.Iterations, pp.Keylen, sha256.New)
	}
	return []byte(fmt.Sprintf("%x", out))
}

func getSetOfCharacters(rules []string) string {
	var setOfChars = ""
	for _, val := range rules {
		v, ok := charRules[val]
		if ok {
			setOfChars += v
		}
	}
	return setOfChars
}

func consumeEntropy(generatedPassword string, quotient *big.Int, setOfCharacters string, maxLength int) (string, *big.Int) {
	if len(generatedPassword) >= maxLength {
		return generatedPassword, quotient
	}
	setLen := big.NewInt(int64(len(setOfCharacters)))
	divR := big.NewInt(0).Mod(quotient, setLen)
	divQ := big.NewInt(0).Div(quotient, setLen)
	generatedPassword += string(setOfCharacters[int(divR.Uint64())])
	return consumeEntropy(generatedPassword, divQ, setOfCharacters, maxLength)
}

func getOneCharPerRule(entropy *big.Int, rules []string) (string, *big.Int) {
	var oneCharPerRules = ""
	for _, rule := range rules {
		password, curEntropy := consumeEntropy("", entropy, charRules[rule], 1)
		oneCharPerRules += password
		entropy = curEntropy
	}
	return oneCharPerRules, entropy
}

func insertStringPseudoRandomly(generatedPassword string, entropy *big.Int, _string string) string {
	for i := 0; i < len(_string); i++ {
		genPasswordBigLen := big.NewInt(int64(len(generatedPassword)))
		divR := big.NewInt(0).Mod(entropy, genPasswordBigLen)
		divQ := big.NewInt(0).Div(entropy, genPasswordBigLen)
		generatedPassword = generatedPassword[0:int(divR.Uint64())] + string(_string[i]) + generatedPassword[int(divR.Uint64()):]
		entropy = divQ
	}
	return generatedPassword
}

func renderPassword(entropy []byte, pp *PasswordProfile) string {
	var setOfCharacters = getSetOfCharacters(pp.Rules)
	_int := big.NewInt(0)
	newInt, _ := _int.SetString(string(entropy), 16)
	genPassword, genEntropyInt := consumeEntropy("", newInt, setOfCharacters, pp.Length-len(pp.Rules))
	addPassword, addEntropyInt := getOneCharPerRule(genEntropyInt, pp.Rules)
	return insertStringPseudoRandomly(genPassword, addEntropyInt, addPassword)
}
