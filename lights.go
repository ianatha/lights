package main

import (
	"fmt"
	"github.com/keep94/gohue"
	"github.com/keep94/maybe"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/urfave/cli"
	"log"
	"os"
)

type Lights struct {
	config Config
	bridge *gohue.Context
	count int
}

func handle_err(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (this Lights) set_color_func(color gohue.Color) func(*cli.Context) {
	return func(c *cli.Context) {
		props := gohue.LightProperties{C: gohue.NewMaybeColor(color), TransitionTime: maybe.NewUint16(0)}
		for i := 1; i <= this.count; i++ {
			this.bridge.Set(i, &props)
		}
	}
}

func (this Lights) set_color_hex_func() func(*cli.Context) {
	return func(c *cli.Context) {
		hexcolor := c.Args().Get(0)
		colour, err := colorful.Hex(hexcolor)
		if err != nil {
			log.Println(hexcolor)
			log.Fatal(err)
		}
		x, y, _ := colour.Xyy()
		props := gohue.LightProperties{C: gohue.NewMaybeColor(gohue.NewColor(x, y)), TransitionTime: maybe.NewUint16(0)}
		for i := 1; i <= this.count; i++ {
			this.bridge.Set(i, &props)
		}
	}
}

func (this Lights) set_color_from_pointer() func(*cli.Context) {
	return func(c *cli.Context) {
		last_colour := ColorAtScreen()
		for {
			colour := ColorAtScreen()
			if colour != last_colour {
				fmt.Printf("color: %#+v\n", colour)
				x, y, Y := colour.Xyy()
				props := gohue.LightProperties{Bri: maybe.NewUint8(uint8(255 * Y)), C: gohue.NewMaybeColor(gohue.NewColor(x, y)), TransitionTime: maybe.NewUint16(0)}
				for i := 1; i <= this.count; i++ {
					this.bridge.Set(i, &props)
				}
			}
			last_colour = colour
		}
	}
}

func MakeLights(cfg Config) Lights {
	return Lights{
		config: cfg,
		bridge: gohue.NewContext(cfg.MeetHue.IPAddress, cfg.MeetHue.Username),
		count: cfg.MeetHue.LightsCount,
	}
}