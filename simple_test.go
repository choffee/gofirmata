package firmata

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestSimple(t *testing.T) {

	board := new(Board)
	board.Device = "/dev/ttyUSB1"
	board.Baud = 57600
	err := board.Setup()
	if err != nil {
		log.Fatal("Could not setup board")
	}
	go func() {
		for {
			msg := <-*board.Reader
			// For now just print out the messages
			fmt.Println(msg)
		}
	}()

	// Set the mode of a pin
	println("set 13 to output")
	board.SetPinMode(13, MODE_OUTPUT)

	// Turn on pin 13
	println("set 13 to 1")
	board.WriteDigital(13, 1)

	// Make it flash
	var onoff byte
	for i := 0; i < 2; i++ {
		board.WriteDigital(13, onoff)
		time.Sleep(1000 * time.Millisecond)
		onoff = (^onoff) & 1
	}

	println("set 13 to analog")
	board.SetPinMode(13, MODE_ANALOG)

	for i := 0; i < 1024; i++ {
		board.WriteAnalog(13, byte(i&0xFF))
		time.Sleep(10 * time.Millisecond)
	}
}
