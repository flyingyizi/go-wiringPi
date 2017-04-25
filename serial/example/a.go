package main

import (
	"log"
	"time"

	"io"

	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: "/dev/ttyS0", Baud: 9600, ReadTimeout: time.Millisecond * 500}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	n, err := s.Write([]byte("test"))
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	for {
		n, err = s.Read(buf)
		if n > 0 {
			log.Print("%:", buf[:n])
		} else if err == io.EOF {
			log.Print("read eof")
			break
		}

	}

	//func (*File) Read
	//
	//Read reads up to len(b) bytes from the File. It returns the number of bytes read and an error, if any. EOF is signaled by a zero count with err set to io.EOF.
	//
	//func (f *File) Read(b []byte) (n int, err error)

	log.Print("%q", buf[:n])
}
