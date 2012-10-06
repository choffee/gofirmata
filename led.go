// Pins and colors
package firmata

import (
	"encoding/hex"
	"fmt"
)

var Colors = map[string][3]byte{
	"red":   [3]byte{0xFF, 00, 00},
	"green": [3]byte{00, 0xFF, 00},
	"blue":  [3]byte{00, 00, 0xFF},
	"black": [3]byte{00, 00, 00},
	"white": [3]byte{0xFF, 0xFF, 0xFF},
}

// An error if we get a bad color
type ColorError struct {
	Desc string
}

func (e ColorError) Error() string {
	return e.Desc
}

// A type for holding some colors
// You can set the color bytes
// The invert bool is for LED's that are wired backwards.
type RGBLED struct {
	rpin, bpin, gpin uint8
	Red, Green, Blue byte
	Invert           bool
}

// Create a new RGB LED
// Expects to have the pin numbers
func NewRGBLED(rp, gp, bp uint8) *RGBLED {
	led := new(RGBLED)
	led.rpin = rp
	led.gpin = gp
	led.bpin = bp
	return led
}

func (l *RGBLED) SetupPins(b *Board) {
	b.SetPinMode(l.rpin, MODE_PWM)
	b.SetPinMode(l.gpin, MODE_PWM)
	b.SetPinMode(l.bpin, MODE_PWM)
}

// Set the LED pins
func (l *RGBLED) Pins(r, g, b uint8) {
	l.rpin, l.bpin, l.gpin = r, g, b
}

// Set the color values from 3 byte array
func (l *RGBLED) Color(c [3]byte) {
	l.Red = c[0]
	l.Green = c[1]
	l.Blue = c[2]
}

// Return the color of the LED as a hex string "RRGGBB"
func (l *RGBLED) HexString() string {
	return fmt.Sprintf("%02X%02X%02X", l.Red, l.Green, l.Blue)
}

func FromHex(s string) ([3]byte, error) {
	var color [3]byte
	var b []byte
	var err error
	for l := 0; l <= 2; l++ {
		b, err = hex.DecodeString(s[0+l*2 : 2+l*2])
		if err != nil {
			fmt.Println(err)
			return color, err
		}
		color[l] = b[0]
	}
	return color, err
}

// Given a string try and convert it to a color
// Strings like "blue", "red" "green" or
// "#FFFE34" or "DEDEDE"
func (l *RGBLED) QuickColor(s string) error {
	var newcolor [3]byte
	var err error
	var ok bool
	if newcolor, ok = Colors[s]; ok {
		l.Color(newcolor)
		return nil
	}
	if newcolor, err = FromHex(s); err == nil {
		l.Color(newcolor)
		return nil
	}
	return ColorError{fmt.Sprintf("Unknow Color %s", s)}
}

// Send the current color of this led to the board
func (l *RGBLED) SendColor(b *Board) {
	l.SetupPins(b)
	var red, green, blue byte
	if l.Invert {
		red, green, blue = 255-l.Red, 255-l.Green, 255-l.Blue
	} else {
		red, green, blue = l.Red, l.Green, l.Blue
	}
	b.WriteAnalog(l.rpin, red)
	b.WriteAnalog(l.gpin, green)
	b.WriteAnalog(l.bpin, blue)
}
