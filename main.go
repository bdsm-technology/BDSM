package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/valyala/fasttemplate"
)

var version = "undefined"

func main() {
	app := cli.NewApp()

	app.Version = version
	app.Name = "bdsm"
	app.HelpName = "bdsm"
	app.Usage = "Bedrock Dedicated Server Manager (Unofficial)"
	app.Author = "CodeHz"
	app.Email = "codehz@outlook.com"

	app.Commands = []cli.Command{
		{
			Name:     "run",
			Aliases:  []string{"r"},
			Usage:    "Run bedrock server directly",
			Category: "For debugging",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "profile, p",
					Value:  "default",
					Usage:  "Profile `name`",
					EnvVar: "profile",
				},
				cli.BoolFlag{
					Name:   "user, u",
					Usage:  "Using user dbus",
					EnvVar: "user_dbus",
				},
				cli.StringFlag{
					Name:  "prompt",
					Value: "{{esc}}[0;36;1mbedrock_server:{{esc}}[22m//{{username}}@{{hostname}}$ {{esc}}[33;0m",
					Usage: "Prompt `template`",
				},
			},
			Action: func(c *cli.Context) error {
				run(c.String("profile"), c.Bool("user"), fasttemplate.New(c.String("prompt"), "{{", "}}"))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
