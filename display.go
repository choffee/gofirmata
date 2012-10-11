// Support for char displays 20x4 etc
package firmata

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

const (
	COMMSI2C      byte = 0
	COMMSSERIAL   byte = 1
	COMMSPARALLEL byte = 2
)

type settings struct {
	name          string
	rows, columns byte
	wrap          bool
}

// Display
type display struct {
	board       *Board
	displayType string
	comms       string
	i2CAddr     byte
	width       byte
	height      byte
	Content     [][]byte
	cursorR     byte
	cursorC     byte
	messages    map[string]byte
	settings    settings
}

// An error if we get a bad color
type BadLocationError struct {
	Desc string
}

func (e BadLocationError) Error() string {
	return e.Desc
}

// Get the settings for a board
func (d *display) getSettings(b string) bool {
	var m map[string]byte
	var s settings
	found := false
	switch b {
	case "LED03":
		s.name = "LED03"
		s.columns = 20
		s.rows = 4
		s.wrap = true
		m["CLEAR"] = 12
		m["MOVE"] = 3
		d.settings = s
		d.messages = m
		found = true
	}
	return found
}

// Create a new I2C display
// Requires the address of the display that you  want to use
func NewI2CDisplay(board *Board, addr byte, dispType string) display {
	disp := newDisplay(dispType)
	disp.comms = "I2C"
	disp.i2CAddr = addr
	return disp
}

func newDisplay(dispType string) display {
	var disp display
	disp.displayType = dispType
	disp.getSettings(dispType)
	disp.setSize()
	return disp
}

func (disp *display) setSize(c, r byte) {
	disp.width = disp.settings.columns
	disp.height = disp.settings.rows
	// Create a new blank copy of the content
	for l := byte(0); l < disp.height; l++ {
		disp.Content = append(disp.Content, make([]byte, disp.width))
	}
	disp.cursorC = 0
	disp.cursorR = 0
}

func (disp *display) send(msg []byte) {
	newmsg := make([]byte, len(msg)+1)
	newmsg[0] = 0
	for l, v := range msg {
		newmsg[l+1] = v
	}
	switch disp.comms {
	case "I2C":
		disp.board.I2CWrite(disp.i2CAddr, I2C_MODE_WRITE, newmsg)
	}
}

func (disp *display) Clear() {
	msg := []byte{disp.messages["CLEAR"]}
	disp.send(msg)
	for rk, _ := range disp.Content {
		for ck, _ := range disp.Content[rk] {
			disp.Content[rk][ck] = 0
		}
	}
	disp.MoveTo(0, 0)
}

// Move the cursor to a location
func (disp *display) MoveTo(r, c byte) error {
	if (r < disp.height) && (c < disp.width) {
		msg := []byte{disp.messages["MOVE"], c, r}
		disp.send(msg)
		disp.cursorR = r
		disp.cursorC = c
	} else {
		return BadLocationError{fmt.Sprintf("row: %d, col: %d, is outside %d, %d", r, c, disp.height, disp.width)}
	}
}

func (disp *display) GetText(r, c, length byte) (string, error) {
	if (r < disp.height) && (c+length < disp.width) {
		var text string
		for l := c; l < length; l++ {
			text = text + string(disp.Content[r][l])
		}
		return text, nil
	} else {
		return "", BadLocationError{fmt.Sprintf("row: %d, col: %d, is outside %d, %d", r, c, disp.height, disp.width)}
	}
}

// Update the display from a new array
func (disp *display) UpdateScreen(newscreen [][]byte) {
	// For now just do it via brute force
	for rk, _ := range newscreen {
		disp.PutText(string(newscreen[rk]), byte(rk), 0)
	}
}

func (disp *display) Write(s string) {
	msg := []byte(s)
	disp.send(msg)
	for _, v := range s {
		disp.Content[disp.cursorR][disp.cursorC] = byte(v)
		if disp.cursorC < disp.width {
			disp.cursorC++
		} else {
			if disp.settings.wrap {
				if disp.cursorR < disp.height {
					disp.cursorR++
				} else {
					disp.cursorR = 0
				}
			} else {
				break
			}
		}
	}
}

// Move to a location and then write the string
func (disp *display) PutText(s string, r, c byte) {
	disp.MoveTo(r, c)
	disp.Write(s)
}

// Keep updating the time on the board
func (disp *display) ShowTime(r, c byte) {
	for {
		now := time.Now()
		clock := now.Format(time.Kitchen)
		disp.PutText("Time "+clock, r, c)
		time.Sleep(1000 * time.Millisecond)
	}
}

func update_bubbles(screen *[][]byte) {
	for rk, _ := range *screen {
		// First replace all the bubbles with the next stage
		(*screen)[rk] = bytes.Replace((*screen)[rk], []byte("*"), []byte(" "), -1)
		(*screen)[rk] = bytes.Replace((*screen)[rk], []byte("Q"), []byte("*"), -1)
		(*screen)[rk] = bytes.Replace((*screen)[rk], []byte("0"), []byte("Q"), -1)
		(*screen)[rk] = bytes.Replace((*screen)[rk], []byte("O"), []byte("0"), -1)
		(*screen)[rk] = bytes.Replace((*screen)[rk], []byte("o"), []byte("O"), -1)
		(*screen)[rk] = bytes.Replace((*screen)[rk], []byte("."), []byte("o"), -1)
		// Now move them up if they can
		if rk > 0 { // not the top row
			for ck, _ := range (*screen)[rk] {
				if bytes.Contains([]byte(".oO0Q*"), []byte{(*screen)[rk][ck]}) {
					(*screen)[rk][ck] = []byte(" ")[0]
					(*screen)[rk-1][ck] = (*screen)[rk][ck]
				}
			}
		}
	}
}

// A better idea may be to add a channel to the function above for doing this.
func add_bubbles(screen *[][]byte) {
	for k, _ := range (*screen)[len(*screen)-1] {
		if rand.Intn(10) > 4 {
			(*screen)[len(*screen)-1][k] = []byte(".")[0]
		}
	}
}
