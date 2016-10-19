package main

import (
	"github.com/lucasb-eyer/go-colorful"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa
#include "foundation.h"
*/
import "C"

func ColorAtScreen() colorful.Color {
	nscolor := C.color_at_screen()
	r := float64(C.color_red_component(nscolor))
	g := float64(C.color_green_component(nscolor))
	b := float64(C.color_blue_component(nscolor))

	return colorful.Color{r, g, b}
}
