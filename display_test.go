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
func (b *BoardTests) TestCreateDisplay(t *C) {
	l := NewI2CDisplay(b, 0x6C, "LED03")
}

func (b *BoardTests) TestMoveTo(t *C) {
	l := NewI2CDisplay(b, 0x6C, "LED03")
	l.MoveTo(5, 3)
}

// We should get an error if we try to move to a place
// outside the display
func (b *BoardTests) TestMoveToBad(t *C) {
	l := NewI2CDisplay(b, 0x6C, "LED03")
	t.Assert(l.MoveTo(25, 2), ErrorMatches, ".*outside.*")
	t.Assert(l.MoveTo(10, 5), ErrorMatches, ".*outside.*")
	t.Assert(l.MoveTo(25, 5), ErrorMatches, ".*outside.*")
}

func (b *BoardTests) TestWriteText(t *C) {
	l := NewI2CDisplay(b, 0x6C, "LED03")
	l.MoveTo(0, 0)
	l.Write("Write Test")
	err, resp := l.GetText(0, 0, 10)
	if err != nil {
		t.Fail()
	} else {
		t.Assert(resp, Matches, "Write Test")
	}
}
