package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/howeyc/gopass"
	CoreV1 "github.com/tuxlinuxien/lesspassgo/core/v1"
	CoreV2 "github.com/tuxlinuxien/lesspassgo/core/v2"
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
			Value: 16,
		},
		cli.BoolFlag{
			Name: "version1, v1",
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
		cli.BoolTFlag{
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
			fmt.Scan(&login)
		}
		if masterPassword == "" {
			fmt.Printf("master password: ")
			masterPasswordBytes, _ := gopass.GetPasswd()
			masterPassword = string(masterPasswordBytes)
		}
		if ctx.Bool("version1") {
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
		} else {
			pp := CoreV2.NewPasswordProfile()
			if ctx.Bool("lower") {
				pp.Rules = append(pp.Rules, "lowercase")
			}
			if ctx.Bool("upper") {
				pp.Rules = append(pp.Rules, "uppercase")
			}
			if ctx.Bool("numbers") {
				pp.Rules = append(pp.Rules, "numbers")
			}
			if ctx.Bool("symbols") {
				pp.Rules = append(pp.Rules, "symbols")
			}
			pp.Length = length
			pp.Counter = counter
			fmt.Println(CoreV2.GeneratePassword(site, login, masterPassword, pp))
		}
		return nil
	}
	app.Run(os.Args)
}
