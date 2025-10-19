package main

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/examples/ili9341/initdisplay"
	"tinygo.org/x/drivers/ili9341"
)

var (
	black = color.RGBA{0, 0, 0, 255}
	white = color.RGBA{255, 255, 255, 255}
	red   = color.RGBA{255, 0, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
	green = color.RGBA{0, 255, 0, 255}
)

func main() {
	button1 := machine.WIO_KEY_C
	button1.Configure(machine.PinConfig{Mode: machine.PinInput})

	button2 := machine.WIO_KEY_B
	button2.Configure(machine.PinConfig{Mode: machine.PinInput})

	button3 := machine.WIO_KEY_A
	button3.Configure(machine.PinConfig{Mode: machine.PinInput})

	display := initdisplay.InitDisplay()

	width, height := display.Size()
	if width < 320 || height < 240 {
		display.SetRotation(ili9341.Rotation270)
	}

	display.FillScreen(black)

	state := 0

	for {
		if !button1.Get() {
			state = 1
		}

		if !button2.Get() {
			state = 2
		}

		if !button3.Get() {
			state = 3
		}

		switch state {
		case 1:
			display.FillScreen(red)
		case 2:
			display.FillScreen(green)
		case 3:
			display.FillScreen(blue)
		default:
			display.FillScreen(black)
		}
	}
}
