package raspi

import (
	"fmt"
	"time"

	serial "go.bug.st/serial"
)

func SendMessage(message string, phoneNumber string) {
	//Serial port opening
	mode := &serial.Mode{
		BaudRate:          115200,
		DataBits:          8,
		Parity:            serial.NoParity,
		StopBits:          serial.OneStopBit,
		InitialStatusBits: &serial.ModemOutputBits{},
	}
	port, err := serial.Open("/dev/ttyS0", mode)
	if err != nil {
		fmt.Println("Error opening serial port: ", err)
		return
	}
	defer port.Close()

	//wait time for module to start
	time.Sleep(20 * time.Second)

	// Enter SMS text mode
	port.Write([]byte("AT+CMGF=1\r\n"))
	time.Sleep(1000 * time.Millisecond)
	buf := make([]byte, 1024)
	n, _ := port.Read(buf)
	fmt.Println("response: ", string(buf[:n]))

	// Set the phone number
	port.Write([]byte(fmt.Sprintf("AT+CMGS=\"%s\"\r\n", phoneNumber)))
	time.Sleep(1000 * time.Millisecond)
	n, _ = port.Read(buf)
	fmt.Println("Response: ", string(buf[:n]))

	// Send the message
	port.Write([]byte(message))
	time.Sleep(1000 * time.Millisecond)
	port.Write([]byte{0x1A})
	time.Sleep(1000 * time.Millisecond)
	n, _ = port.Read(buf)
	fmt.Println("Response: ", string(buf[:n]))

}
