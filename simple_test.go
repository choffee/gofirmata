package firmata

import (
	"log"
	"testing"
	"time"
)

func getBoard() *Board {

	board := new(Board)
	board.Device = "/dev/ttyUSB1"
	board.Baud = 57600
	err := board.Setup()
	if err != nil {
		log.Fatal("Could not setup board")
	}
	return board

}

func TestAnalogWrite(t *testing.T) {
	board := getBoard()
	println("set 13 to analog")
	board.SetPinMode(13, MODE_ANALOG)

	println("Analog pulse on pin 13")
	for i := 0; i < 1024; i++ {
		board.WriteAnalog(13, byte(i&0xFF))
		time.Sleep(10 * time.Millisecond)
	}
}

func TestDigitalWrite(t *testing.T) {
	board := getBoard()
	// Set the mode of a pin
	println("set 13 to output")
	board.SetPinMode(13, MODE_OUTPUT)

	// Turn on pin 13
	println("set 13 to 1")
	board.WriteDigital(13, 1)

	// Make it flash
	println("Flash pin 13")
	var onoff byte
	for i := 0; i < 2; i++ {
		board.WriteDigital(13, onoff)
		time.Sleep(1000 * time.Millisecond)
		onoff = (^onoff) & 1
	}
}

func TestI2CConfig(t *testing.T) {
	board := getBoard()
	println("Setting up I2C")
	board.I2CConfig(0)
}

func TestI2CSend(t *testing.T) {
	board := getBoard()
	println("Sending I2C clear screen")
	LCDaddr := byte(0xC6 >> 1) // For the LCD02 screen that I have
	msg := []byte{12}          // Clear the screen
	board.I2CWrite(LCDaddr, msg)
}
