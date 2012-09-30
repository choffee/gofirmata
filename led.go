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

// A type for holding some colors
type RGBLED struct {
	rpin, bpin, gpin uint8
	Red, Green, Blue byte
}

// An error if we get a bad color
type ColorError struct {
	Desc string
}

func (e ColorError) Error() string {
	return e.Desc
}

// Set the LED pins
func (l *RGBLED) Pins(r, g, b uint8) {
	l.rpin, l.bpin, l.gpin = r, g, b
}

func (l *RGBLED) Color(c [3]byte) {
	l.rpin = c[0]
	l.gpin = c[1]
	l.bpin = c[2]
}

func (l *RGBLED) QuickColor(s string) error {
	var newcolor [3]byte
	var ok bool
	if newcolor, ok = Colors[s]; ok {
		l.Color(newcolor)
		return nil
	}
	return ColorError{fmt.Sprintf("Unknow Color %s", s)}
}
