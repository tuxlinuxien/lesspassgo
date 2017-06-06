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
	app.UsageText = "lesspassgo [options]"
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
		cli.IntFlag{
			Name:   "port,p",
			Value:  1314,
			EnvVar: "PORT",
		},
	}

	app.Action = func(ctx *cli.Context) error {
		server.Start(ctx.String("db"), ctx.Int("port"))
		return nil
	}
	app.Run(os.Args)
}
