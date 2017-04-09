package rpi_test

import (
	"fmt"
	"time"

	_ "github.com/flyingyizi/go-wiringPi/rpi/raspberry.go"
)

// LED Pin - BCM_GPIO 17.
const LED int = 17

func Example_blink() {
	fmt.Println("Raspberry Pi blink")
	rpi.Init()
	defer rpi.Close()
	pin := rpi.Pin(4)
	pin.Output()
	for {
		pin.Toggle()
		time.Sleep(time.Second)
	}

	//// Output: Raspberry Pi blink
}
