// Pins and colors
package firmata

import (
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
type RGBLED struct {
	rpin, bpin, gpin uint8
	Red, Green, Blue byte
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
	l.rpin = c[0]
	l.gpin = c[1]
	l.bpin = c[2]
}

// Return the color of the LED as a hex string "RRGGBB"
func (l *RGBLED) HexString() string {
	return fmt.Sprintf("%02X%02X%02X", l.rpin, l.gpin, l.bpin)
}

// Given a string try and convert it to a color
// Strings like "blue", "red" "green" or
// "#FFFE34" or "DEDEDE"
func (l *RGBLED) QuickColor(s string) error {
	var newcolor [3]byte
	var ok bool
	if newcolor, ok = Colors[s]; ok {
		l.Color(newcolor)
		return nil
	}
	return ColorError{fmt.Sprintf("Unknow Color %s", s)}
}

// Send the current color of this led to the board
func (l *RGBLED) SendColor(b *Board) {
	b.WriteAnalog(l.rpin, l.Red)
	b.WriteAnalog(l.gpin, l.Green)
	b.WriteAnalog(l.bpin, l.Blue)
}
