package firmata

import (
  "github.com/tarm/goserial"
  "log"
  "io"
  "fmt"
  "strconv"
)

const (
  MODE_INPUT byte  = 0x00
  MODE_OUTPUT byte = 0x01
  MODE_ANALOG byte = 0x02
  MODE_PWM byte    = 0x03
  MODE_SERVO byte  = 0x04
  MODE_SHIFT byte  = 0x05
  MODE_I2C byte    = 0x06

  UNKNOWN byte                 = 0xFF // I just invented this it could be used elsewhere
  START_SYSEX byte             = 0xF0 // start a MIDI Sysex message
  END_SYSEX byte               = 0xF7 // end a MIDI Sysex message
  PIN_MODE_QUERY byte          = 0x72 // ask for current and supported pin modes
  PIN_MODE_RESPONSE byte       = 0x73 // reply with current and supported pin modes
  PIN_STATE_QUERY byte         = 0x6D
  PIN_STATE_RESPONSE byte      = 0x6E
  CAPABILITY_QUERY byte        = 0x6B
  CAPABILITY_RESPONSE byte     = 0x6C
  ANALOG_MAPPING_QUERY byte    = 0x69
  ANALOG_MAPPING_RESPONSE byte = 0x6A
  REPORT_FIRMWARE byte         = 0x79 // report name and version of the firmware
  PIN_MODE byte                = 0xF4 // Set the pin mode

  DIGITAL_WRITE byte           = 0x90

)

type FirmataMsg struct {
  msgtype byte
  pin byte
  data map[string]string
  rawdata []byte
}

type Board struct {
  Name string
  config *serial.Config
  Device string
  Baud int
  serial io.ReadWriteCloser
  Reader *chan FirmataMsg
  Writer *chan FirmataMsg
  digitalPins [8]byte // Keeps a record of digital pin values
}

// Setup the board to start reading and writing
// I expect you to have already setup a serial config
func (board *Board) Setup () error {
  board.config = &serial.Config{Name: board.Device, Baud: board.Baud}
  var err error
  board.serial, err = serial.OpenPort(board.config)
  if(err != nil){ log.Fatal("Could not open port") }
  board.GetReader()
  return err
}

func process_sysex(sysextype byte, msgdata []byte ) FirmataMsg {
  var result FirmataMsg
  fmt.Println("SYSEX: %d", sysextype, msgdata)
  switch sysextype {
  case REPORT_FIRMWARE: // queryFirmware
    result.msgtype = REPORT_FIRMWARE
    result.data = make(map[string]string)
    result.data["major"] = strconv.Itoa(int(msgdata[1]))
    result.data["minor"] = strconv.Itoa(int(msgdata[2]))
    result.data["name"]  = string(msgdata[3:]) //TODO I don't think this works
  default:
    result.msgtype = UNKNOWN
    result.data = make(map[string]string)
    result.data["msgtyperaw"] = string(sysextype)
    result.data["unknown"] = string(msgdata)
  }
  return result
}

// Pass a pointer to a serial port and this function will send back the messages
// received over a chanel
func (board *Board) GetReader()  {
  board.Reader = new(chan FirmataMsg)
  go func() {
    for {
      l  := make([]byte,1)
      _, err := board.serial.Read(l)
      if err != nil {
        log.Fatal("Failed to read from Serial port")
        return
      } else {
        switch l[0] {
        case  START_SYSEX:
          t := make([]byte,1)
          var sysextype byte
          _, terr := board.serial.Read(t)
          if terr != nil {
            log.Fatal("Failed to read sysex type")
          } else {
            sysextype = t[0]
          }
          var merr error
          var msgdata []byte
          for m := make([]byte, 1) ; m[0] != END_SYSEX ; _, merr = board.serial.Read(m) {
            if merr != nil {
              log.Fatal("Failed to read sysex from serial port")
            } else {
              msgdata = append(msgdata, m[0])
            }
          }
          // Send the message down the chanel
          newmsg := process_sysex(sysextype, msgdata)
          *board.Reader <- newmsg
        }
      }
    }
  }()
}

func (board *Board) sendMsg (msg FirmataMsg) {
}

func (board *Board) sendRaw (msg *[]byte) {
  board.serial.Write(*msg)
}

// Set the mode for a pin
// mode should be one of: MODE_INPUT MODE_OUTPUT, MODE_ANALOG,
//                        MODE_PWM, MODE_SERVO, MODE_SHIFT, MODE_I2C
func (board *Board) SetPinMode(pin, mode byte) {
  msg := new(FirmataMsg)
  msg.msgtype = PIN_MODE
  msg.pin     = pin
  msg.rawdata = []byte{mode}
  board.sendMsg(*msg)
}

func (board *Board) WriteDigital(pin, value byte) {
  port := (pin >> 3) & 0x0F // Get the port the pin is in
  // Next we need to get all 8 pins for that port and only change the one
  // we are intrested in
  switch value {
  case 0:
    board.digitalPins[port] = board.digitalPins[port] & ^(1 << (pin & 0x07)) 
  case 1:
    board.digitalPins[port] = (board.digitalPins[port] | (1 << (pin & 0x07)))
  }
  // Now send the whole port ( 8 pins ) to the arduino
  cmd := byte(DIGITAL_WRITE | port)
  msg:= []byte{cmd, board.digitalPins[port]  & 0x7F, (board.digitalPins[port] >> 7)  & 0x7f }
  board.sendRaw(&msg)
}


