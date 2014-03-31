package main

import (
	"code.google.com/p/gcfg"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/keep94/gohue"
	"github.com/keep94/maybe"
	"os"
)

type Config struct {
	MeetHue struct {
		IPAddress   string
		Username    string
		Flag        bool
		LightsCount int
	}
}

func main() {
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, ".lights.ini")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app := cli.NewApp()
	app.Name = "lights"
	app.Usage = "make an explosive entrance"

	bridge := gohue.NewContext(cfg.MeetHue.IPAddress, cfg.MeetHue.Username)

	lights_count := cfg.MeetHue.LightsCount

	var red = cli.Command{
		Name:      "red",
		ShortName: "r",
		Usage:     "set all lights to red",
		Action: func(c *cli.Context) {
			props := gohue.LightProperties{C: gohue.NewMaybeColor(gohue.Red)}
			for i := 1; i <= lights_count; i++ {
				bridge.Set(i, &props)
			}
		},
	}

	var white = cli.Command{
		Name:      "white",
		ShortName: "w",
		Usage:     "set all lights to white",
		Action: func(c *cli.Context) {
			props := gohue.LightProperties{C: gohue.NewMaybeColor(gohue.White)}
			for i := 1; i <= lights_count; i++ {
				bridge.Set(i, &props)
			}
		},
	}

	var command_on = cli.Command{
		Name:      "off",
		ShortName: "0",
		Usage:     "turn off all lights",
		Action: func(c *cli.Context) {
			props := gohue.LightProperties{On: maybe.NewBool(false)}
			for i := 1; i <= lights_count; i++ {
				bridge.Set(i, &props)
			}
		},
	}

	var command_off = cli.Command{
		Name:      "on",
		ShortName: "1",
		Usage:     "turn on all lights",
		Action: func(c *cli.Context) {
			props := gohue.LightProperties{On: maybe.NewBool(true)}
			for i := 1; i <= lights_count; i++ {
				bridge.Set(i, &props)
			}
		},
	}

	app.Commands = []cli.Command{
		command_on,
		command_off,
		red,
		white,
	}

	app.Run(os.Args)
}
