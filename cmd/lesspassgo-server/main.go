package main

import (
	"os"

	"github.com/tuxlinuxien/lesspassgo/server"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "lesspassgo-server"
	app.Usage = "LessPass password server."
	app.UsageText = "lesspassgo-server [options]"
	app.HideVersion = true
	app.Author = "Yoann Cerda"
	app.Email = "tuxlinuxien@gmail.com"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "db",
			Value:  "./user.db",
			EnvVar: "DB_PATH",
		},
		cli.StringFlag{
			Name:   "host",
			Value:  "0.0.0.0",
			EnvVar: "HOST",
		},
		cli.IntFlag{
			Name:   "port,p",
			Value:  1314,
			EnvVar: "PORT",
		},
		cli.BoolFlag{
			Name:   "disable-registration",
			EnvVar: "DISABLE_REGISTRATION",
		},
	}

	app.Action = func(ctx *cli.Context) error {
		server.Start(
			ctx.String("db"),
			ctx.String("host"),
			ctx.Int("port"),
			ctx.Bool("disable-registration"),
		)
		return nil
	}
	app.Run(os.Args)
}
