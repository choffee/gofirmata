package main
import (
  "github.com/tarm/goserial"
  "log"
  "time"
)

const (
  MODE_INPUT byte  = 0x00
  MODE_OUTPUT byte = 0x01
  MODE_ANALOG byte = 0x02
  MODE_PWM byte    = 0x03
  MODE_SERVO byte  = 0x04
  MODE_SHIFT byte  = 0x05
  MODE_I2C byte    = 0x06

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

)

func main() {
  c := &serial.Config{Name: "/dev/ttyUSB1", Baud: 57600}
  s, err := serial.OpenPort(c)
  if(err != nil){ log.Fatal("Could not open port") }

  reportfirmware := []byte{START_SYSEX, REPORT_FIRMWARE, END_SYSEX}
  n, err := s.Write(reportfirmware)
  if err != nil {
          log.Fatal("write err")
  }

  buf := make([]byte, 1024)
  n, err = s.Read(buf)
  if err != nil {
    log.Fatal(err)
  }
  log.Print("%q", buf[:n])

  // pin 13

  // Set the mode of a pin
  msg := []byte{0xF4, 13,MODE_OUTPUT}
  n, err = s.Write(msg)
  if err != nil {
    log.Fatal("failed to set pin mode")
  }


  port_num := byte(13 / 8)
  port_value := byte(255)
  for {
    msg = []byte{0x90 | port_num ,port_value & 0x7F, (port_value >> 7) & 0x7f }
    n, err = s.Write(msg)
    if err != nil {
      log.Fatal("failed to set pin mode")
    }
    time.Sleep(1000 * time.Millisecond)
    port_value = 255 - port_value
  }
}
