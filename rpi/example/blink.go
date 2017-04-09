package main

import (
	"fmt"
	"time"

	"github.com/flyingyizi/go-wiringPi/rpi"
)

//go build -gcflags "-N -l"  blink.go

// LED Pin - BCM_GPIO 17.
const LED int = 17

func main() {
	fmt.Println("Raspberry Pi blink")

	pcbrev, bmodel, processor, manufacturer, ram, bWarranty, err := rpi.Boardinfo()
	fmt.Println( pcbrev, bmodel, processor, manufacturer, ram, bWarranty)
	err = rpi.Init()
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
