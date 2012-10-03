package firmata

import (
	"fmt"
	. "launchpad.net/gocheck"
	"log"
	"testing"
	"time"
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

func (b *BoardTests) TestAnalogMapping(t *C) {
	println("Sending analog mapping request")
	b.board.GetAnalogMapping()
	time.Sleep(1000 * time.Millisecond)
}

func (b *BoardTests) TestAnalogWrite(t *C) {
	println("set 13 to analog")
	b.board.SetPinMode(13, MODE_ANALOG)

	println("Analog pulse on pin 13")
	for i := 0; i < 1024; i++ {
		b.board.WriteAnalog(13, byte(i&0xFF))
		time.Sleep(10 * time.Millisecond)
	}
}

func (b *BoardTests) TestDigitalWrite(t *C) {
	// Set the mode of a pin
	println("set 13 to output")
	b.board.SetPinMode(13, MODE_OUTPUT)

	// Turn on pin 13
	println("set 13 to 1")
	b.board.WriteDigital(13, 1)

	// Make it flash
	println("Flash pin 13")
	var onoff byte
	for i := 0; i < 2; i++ {
		b.board.WriteDigital(13, onoff)
		time.Sleep(1000 * time.Millisecond)
		onoff = (^onoff) & 1
	}
}

func (b *BoardTests) TestI2CConfig(t *C) {
	println("Setting up I2C")
	b.board.I2CConfig(0)
}

func (b *BoardTests) TestI2CSend(t *C) {
	b.board.Debug = 1
	println("Sending I2C clear screen")
	LCDaddr := byte(0xC6 >> 1) // For the LCD02 screen that I have
	msg := []byte{12}          // Clear the screen
	b.board.I2CWrite(LCDaddr, I2C_MODE_WRITE, msg)
}
