package main

import (
	"github.com/keep94/gohue"
	"github.com/keep94/maybe"
	"github.com/urfave/cli"
	"gopkg.in/gcfg.v1"
	"math/rand"
	"os"
	"os/user"
	"strconv"
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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	usr, err := user.Current()
	handle_err(err)

	var cfg Config
	err = gcfg.ReadFileInto(&cfg, usr.HomeDir+"/.lights.ini")
	handle_err(err)

	lights := MakeLights(cfg)

	app := cli.NewApp()
	app.Name = "lights"
	app.Usage = "make an explosive entrance"


	var hex = cli.Command{
		Name:      "hex",
		ShortName: "#",
		Usage:     "set all lights to given hex",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "color"},
		},
		Action: lights.set_color_hex_func(),
	}

	var red = cli.Command{
		Name:      "red",
		ShortName: "r",
		Usage:     "set all lights to red",
		Action:    lights.set_color_func(gohue.Red),
	}

	var blue = cli.Command{
		Name:      "blue",
		ShortName: "b",
		Usage:     "set all lights to blue",
		Action:    lights.set_color_func(gohue.Blue),
	}

	var white = cli.Command{
		Name:      "white",
		ShortName: "w",
		Usage:     "set all lights to white",
		Action:    lights.set_color_func(gohue.White),
	}

	var random = cli.Command{
		Name:      "random",
		ShortName: "rnd",
		Usage:     "random color",
		Action: func(c *cli.Context) {
			x := rand.Float64()
			y := rand.Float64()
			lights.set_color_func(gohue.NewColor(x, y))(c)
		},
	}

	var popo = cli.Command{
		Name:      "police",
		ShortName: "popo",
		Usage:     "popo scheme",
		Action: func(c *cli.Context) {
			for {
				lights.set_color_func(gohue.Red)(c)
				time.Sleep(250 * time.Millisecond)
				lights.set_color_func(gohue.Blue)(c)
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
			for i := 1; i <= lights.count; i++ {
				lights.bridge.Set(i, &props)
			}
		},
	}

	var command_off = cli.Command{
		Name:      "on",
		ShortName: "1",
		Usage:     "turn on all lights",
		Action: func(c *cli.Context) {
			props := gohue.LightProperties{On: maybe.NewBool(true)}
			for i := 1; i <= lights.count; i++ {
				lights.bridge.Set(i, &props)
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
			for i := 1; i <= lights.count; i++ {
				lights.bridge.Set(i, &props)
			}
		},
	}

	var command_pointer = cli.Command{
		Name:      "mouse",
		ShortName: "m",
		Usage:     "set light color from pointer",
		Action:    lights.set_color_from_pointer(),
	}

	app.Commands = []cli.Command{
		hex,
		command_on,
		command_off,
		command_brightness,
		red,
		white,
		blue,
		random,
		popo,
		command_pointer,
	}

	app.Run(os.Args)
}
