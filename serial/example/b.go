package main

import (
	"bufio"
	"fmt"

	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	s, err := serial.OpenPort(c)

	if err != nil {
		fmt.Println(err)
	}

	_, err = s.Write([]byte("\x16\x02N0C0 G A\x03\x0d\x0a"))

	if err != nil {
		fmt.Println(err)
	}

	//Which will read data until the first \x0a.
	reader := bufio.NewReader(s)
	reply, err := reader.ReadBytes('\x0a')
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)

	s.Close()
}
