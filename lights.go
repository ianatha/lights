package main

import (
	"code.google.com/p/gcfg"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/keep94/gohue"
	"github.com/keep94/maybe"
	"os"
	"os/user"
	"strconv"
	"math/rand"
	"time"
)

type Config struct {
	MeetHue struct {
		IPAddress   string
		Username    string
		Flag        bool
		LightsCount int
	}
}

func handle_err(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func set_color_func(color gohue.Color, lights_count int, bridge *gohue.Context) func(*cli.Context) {
	return func(c *cli.Context) {
			props := gohue.LightProperties{C: gohue.NewMaybeColor(color), TransitionTime: maybe.NewUint16(0)}
			for i := 1; i <= lights_count; i++ {
				bridge.Set(i, &props)
			}
		}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) 

	usr, err := user.Current()
	handle_err(err)

	var cfg Config
	err = gcfg.ReadFileInto(&cfg, usr.HomeDir+"/.lights.ini")
	handle_err(err)

	app := cli.NewApp()
	app.Name = "lights"
	app.Usage = "make an explosive entrance"

	bridge := gohue.NewContext(cfg.MeetHue.IPAddress, cfg.MeetHue.Username)

	lights_count := cfg.MeetHue.LightsCount

	var red = cli.Command{
		Name:      "red",
		ShortName: "r",
		Usage:     "set all lights to red",
		Action: set_color_func(gohue.Red, lights_count, bridge),
	}

	var blue = cli.Command{
		Name:      "blue",
		ShortName: "b",
		Usage:     "set all lights to blue",
		Action: set_color_func(gohue.Blue, lights_count, bridge),
	}

	var white = cli.Command{
		Name:      "white",
		ShortName: "w",
		Usage:     "set all lights to white",
		Action: set_color_func(gohue.White, lights_count, bridge),
	}

	var random = cli.Command{
		Name:      "random",
		ShortName: "rnd",
		Usage:     "random color",
		Action: func(c *cli.Context) {
			x := rand.Float64()
			y := rand.Float64()
			set_color_func(gohue.NewColor(x, y), lights_count, bridge)(c)
		},
	}

	var popo = cli.Command{
		Name:      "police",
		ShortName: "popo",
		Usage:     "popo scheme",
		Action: func(c *cli.Context) {
			for {
				red.Action(c)
				time.Sleep(250 * time.Millisecond)
				blue.Action(c)
				time.Sleep(250 * time.Millisecond)
			}
		},
	}

	var command_on = cli.Command{
		Name:      "off",
		ShortName: "0",
		Usage:     "turn off all lights",
		Action: func(c *cli.Context) {
			props := gohue.LightProperties{On: maybe.NewBool(false), TransitionTime: maybe.NewUint16(0)}
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

	var command_brightness = cli.Command{
		Name:      "brightness",
		ShortName: "bri",
		Usage:     "turn on all lights",
		Action: func(c *cli.Context) {
			bright, _ := strconv.ParseUint(c.Args().First(), 10, 8)
			props := gohue.LightProperties{Bri: maybe.NewUint8(uint8(bright)), TransitionTime: maybe.NewUint16(0)}
			for i := 1; i <= lights_count; i++ {
				bridge.Set(i, &props)
			}
		},
	}

	app.Commands = []cli.Command{
		command_on,
		command_off,
		command_brightness,
		red,
		white,
		blue,
		random,
		popo,
	}

	app.Run(os.Args)
}
