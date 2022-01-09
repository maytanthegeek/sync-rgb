package main

import (
	"log"
	"os"

	"github.com/maytanthegeek/sync-rgb/internal/home"
	"github.com/urfave/cli/v2"
)

var hd *home.HomeDevices

func main() {

	var config string
	var toggleswitch bool
	var color string

	app := &cli.App{
		Name:  "sync-rgb",
		Usage: "fight the loneliness!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "config file containing devices",
				Aliases:     []string{"c"},
				Required:    true,
				Destination: &config,
			},
			&cli.BoolFlag{
				Name:        "switch",
				Usage:       "toggle switch the lights",
				Aliases:     []string{"t"},
				Destination: &toggleswitch,
			},
			&cli.StringFlag{
				Name:        "color",
				Usage:       "list of RGB values of color to set",
				Aliases:     []string{"r"},
				Destination: &color,
			},
		},
		Action: func(c *cli.Context) error {
			hd = home.NewHomeDevices(config)
			if toggleswitch {
				toggleSwitch()
			}
			if color != "" {
				changeColor(color)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
