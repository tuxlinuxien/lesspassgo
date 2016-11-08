package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/tuxlinuxien/lesspassgo/core"
	"github.com/urfave/cli"
)

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

	var args = make([]string, 0)
	app.Action = func(ctx *cli.Context) error {
		site := args[0]
		login := args[1]
		masterPassword := args[2]
		counter := ctx.Int("counter")
		length := ctx.Int("length")
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
		encLogin := core.EncryptLogin(login, masterPassword)
		fmt.Println(core.RenderPassword(encLogin, site, length, counter, template))
		return nil
	}
	if len(os.Args) < 4 {
		app.Run([]string{"", "-h"})
		return
	}
	for i := 1; i < 4 && i < len(os.Args); i++ {
		args = append(args, os.Args[i])
	}
	var flags = []string{"lesspassgo"}
	for i := 4; i < len(os.Args); i++ {
		flags = append(flags, os.Args[i])
	}
	app.Run(flags)
}
