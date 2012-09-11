package main
import (
  "github.com/tarm/goserial"
  "log"
  "time"
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

type firmata_msg struct {
  msgtype string
  pin int
  data map[string]string
}

func process_sysex(sysextype byte, msgdata []byte ) firmata_msg {
  var result firmata_msg
  fmt.Println("SYSEX: %d", sysextype, msgdata)
  switch sysextype {
  case REPORT_FIRMWARE: // queryFirmware
    result.msgtype = "REPORT_FIRMWARE"
    result.data = make(map[string]string)
    result.data["major"] = strconv.Itoa(int(msgdata[1]))
    result.data["minor"] = strconv.Itoa(int(msgdata[2]))
    result.data["name"]  = string(msgdata[3:]) //TODO I don't think this works
  default:
    result.msgtype = "UNKONWN"
    result.data = make(map[string]string)
    result.data["msgtyperaw"] = string(sysextype)
    result.data["unknown"] = string(msgdata)
  }
  return result
}

// Pass a pointer to a serial port and this function will send back the messages
// received over a chanel
func read_serial( s io.ReadWriteCloser ) ( *chan firmata_msg) {
  results_c := make(chan firmata_msg)
  go func() {
    for {
      l  := make([]byte,1)
      _, err := s.Read(l)
      if err != nil {
        log.Fatal("Failed to read from Serial port")
        return
      } else {
        switch l[0] {
        case  START_SYSEX:
          t := make([]byte,1)
          var sysextype byte
          _, terr := s.Read(t)
          if terr != nil {
            log.Fatal("Failed to read sysex type")
          } else {
            sysextype = t[0]
          }
          var merr error
          var msgdata []byte
          for m := make([]byte, 1) ; m[0] != END_SYSEX ; _, merr = s.Read(m) {
            if merr != nil {
              log.Fatal("Failed to read sysex from serial port")
            } else {
              msgdata = append(msgdata, m[0])
            }
          }
          // Send the message down the chanel
          newmsg := process_sysex(sysextype, msgdata)
          results_c <- newmsg
        }
      }
    }
  }()
  return &results_c
}

func main() {
  config := &serial.Config{Name: "/dev/ttyUSB1", Baud: 57600}
  s, err := serial.OpenPort(config)
  if(err != nil){ log.Fatal("Could not open port") }

  c := read_serial(s)
  fmt.Println(c)
  go func() {
    for {
      msg := <- *c
      // For now just print out the messages
      fmt.Println( msg )
    }
  }()
  // *c <- *new(firmata_msg)

  reportfirmware := []byte{START_SYSEX, REPORT_FIRMWARE, END_SYSEX}
  _, err = s.Write(reportfirmware)
  if err != nil {
          log.Fatal("write err")
  }

  //buf := make([]byte, 1024)
  //n, err = s.Read(buf)
  //if err != nil {
  //  log.Fatal(err)
  //}
  //log.Print("%q", buf[:n])

  // pin 13

  // Set the mode of a pin
  msg := []byte{0xF4, 13,MODE_OUTPUT}
  _, err = s.Write(msg)
  if err != nil {
    log.Fatal("failed to set pin mode")
  }


  port_num := byte(13 / 8)
  port_value := byte(255)
  for {
    msg = []byte{0x90 | port_num ,port_value & 0x7F, (port_value >> 7) & 0x7f }
    _, err = s.Write(msg)
    if err != nil {
      log.Fatal("failed to set pin mode")
    }
    time.Sleep(1000 * time.Millisecond)
    port_value = 255 - port_value
  }
}

