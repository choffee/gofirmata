package firmata

import (
	"fmt"
	. "launchpad.net/gocheck"
	"log"
	"testing"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type BoardTests struct {
	board Board
}

var _ = Suite(&BoardTests{})

func (b *BoardTests) SetUpTest(c *C) {
	board, err := NewBoard("/dev/ttyUSB1", 57600)
	b.board = *board
	if err != nil {
		log.Fatal("Could not setup board")
	}
	go func() {
		for msg := range *board.Reader {
			fmt.Println(msg)
		}
	}()
}

// Just create a new LED and setup the pins
func (b *BoardTests) TestCreateLED(t *C) {
	println("Sending analog mapping request")
	l := new(RGBLED, 9, 10, 11)
	l.SetupPins(b.board)
}

// Check that when we use quick colors they come out as we would expect
func (b *BoardTest) TestQuickColors(t *C) {
	l := new(RGBLED, 9, 10, 11)
	var cs = map[string][3]byte{
		"red":   [3]byte{0xFF, 00, 00},
		"green": [3]byte{00, 0xFF, 00},
		"blue":  [3]byte{00, 00, 0xFF},
		"black": [3]byte{00, 00, 00},
		"white": [3]byte{0xFF, 0xFF, 0xFF},
	}
	for c, v := range cs {
		l.QuickColor(c)
		t.Check(l.Red, Equals, v[0])
		t.Check(l.Green, Equals, v[1])
		t.Check(l.Blue, Equals, v[2])
	}
}

// Check that the conversion to a hex string RRGGBB works
// as expected
func (b *BoardTest) TestHexColors(t *C) {
	l := new(RGBLED, 9, 10, 11)
	l.Red = 02    // Test padding
	l.Green = 255 // largest
	l.Blue = 0    // Zero
	t.check(l.HexString(), Equals, "02FF00")
}
