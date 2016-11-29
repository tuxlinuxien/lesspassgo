package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/howeyc/gopass"
	CoreV1 "github.com/tuxlinuxien/lesspassgo/core/v1"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "lesspassgo"
	app.Usage = "LessPass password generator CLI."
	app.UsageText = "lesspassgo [options]"
	app.HideVersion = true
	app.Author = "Yoann Cerda"
	app.Email = "tuxlinuxien@gmail.com"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "site",
		},
		cli.StringFlag{
			Name: "login",
		},
		cli.StringFlag{
			Name: "password",
		},
		cli.IntFlag{
			Name:  "counter, c",
			Value: 1,
		},
		cli.IntFlag{
			Name:  "length, L",
			Value: 12,
		},
		cli.BoolTFlag{
			Name: "upper, u",
		},
		cli.BoolTFlag{
			Name: "lower, l",
		},
		cli.BoolTFlag{
			Name: "numbers, n",
		},
		cli.BoolFlag{
			Name: "symbols, s",
		},
	}

	app.Action = func(ctx *cli.Context) error {
		site := ctx.String("site")
		login := ctx.String("login")
		masterPassword := ctx.String("password")
		counter := ctx.Int("counter")
		length := ctx.Int("length")
		if site == "" {
			fmt.Printf("site: ")
			fmt.Scan(&site)
		}
		if login == "" {
			fmt.Printf("login: ")
			fmt.Scan(&site)
		}
		if masterPassword == "" {
			fmt.Printf("master password: ")
			masterPasswordBytes, _ := gopass.GetPasswd()
			masterPassword = string(masterPasswordBytes)
		}
		var template = ""
		if ctx.Bool("lower") {
			template += "vc"
		}
		if ctx.Bool("upper") {
			template += "VC"
		}
		if ctx.Bool("numbers") {
			template += "n"
		}
		if ctx.Bool("symbols") {
			template += "s"
		}
		if template == "" {
			return errors.New("At least one of -l -u -n -s required.")
		}
		encLogin := CoreV1.EncryptLogin(login, masterPassword)
		fmt.Println(CoreV1.RenderPassword(encLogin, site, length, counter, template))
		return nil
	}
	app.Run(os.Args)
}
