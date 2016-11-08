package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli"
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
	app := cli.NewApp()
	app.Name = "lesspassgo"
	app.Usage = "LessPass password generator CLI."
	app.UsageText = "lesspassgo <site> <login> <masterPassword> [options]"
	app.HideVersion = true
	app.Author = "Yoann Cerda"
	app.Email = "tuxlinuxien@gmail.com"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.Int64Flag{
			Name:  "counter, c",
			Value: 1,
		},
		cli.Int64Flag{
			Name:  "length, L",
			Value: 12,
		},
		cli.BoolFlag{
			Name: "upper, u",
		},
		cli.BoolFlag{
			Name: "lower, l",
		},
		cli.BoolFlag{
			Name: "numbers, n",
		},
		cli.BoolFlag{
			Name: "symbols, s",
		},
	}

	app.Action = func(ctx *cli.Context) error {
		fmt.Println(ctx.NArg())
		return nil
	}
	app.Run(os.Args)

	// login := flag.String("login", "", "login")
	// masterPassword := flag.String("password", "", "password")
	// site := flag.String("site", "", "domain.com")
	// length := flag.Int("length", 12, "generated password length")
	// counter := flag.Int("counter", 1, "counter")
	// lower := flag.Bool("l", false, "password contains [a-z]")
	// upper := flag.Bool("u", false, "password contains [A-Z]")
	// numeric := flag.Bool("n", false, "password contains [0-9]")
	// special := flag.Bool("s", false, "password contains @&%?,=[]_:-+*$#!'^~;()/.")
	// flag.Parse()
	//
	// var template = ""
	// if *lower == true {
	// 	template += "vc"
	// }
	// if *upper == true {
	// 	template += "VC"
	// }
	// if *numeric == true {
	// 	template += "n"
	// }
	// if *special == true {
	// 	template += "s"
	// }
	// if template == "" {
	// 	fmt.Println("You need to define a password format")
	// 	os.Exit(-1)
	// }
	// encLogin := encryptLogin(*login, *masterPassword)
	// fmt.Println(renderPassword(encLogin, *site, *length, *counter, template))
}
