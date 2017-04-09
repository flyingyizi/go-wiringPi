package main

import (
	"fmt"
	"time"

	"github.com/flyingyizi/go-wiringPi/rpi"
)

// LED Pin - BCM_GPIO 17.
const LED int = 17

func main() {
	fmt.Println("Raspberry Pi blink")
	err := rpi.Init()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rpi.Close()
	pin := rpi.Pin(4)
	pin.Output()
	for {
		pin.TogglePin()
		time.Sleep(time.Second)
	}

	// Output: Raspberry Pi blink
}
