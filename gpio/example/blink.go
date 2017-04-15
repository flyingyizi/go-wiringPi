package main

import (
	"fmt"
	"time"

	"github.com/flyingyizi/go-wiringPi/gpio"
)

//go build -gcflags "-N -l"  blink.go

// LED Pin - BCM_GPIO 17.
const LED int = 17

func main() {
	fmt.Println("Raspberry Pi blink")
	info, base, err := gpio.GetBoardInfo()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("modelName:", info.ModelName())

	err = gpio.Init()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer gpio.Close()
	pin := gpio.Pin(4)
	pin.Output()
	for {
		pin.TogglePin()
		time.Sleep(time.Second)
	}

	// Output: Raspberry Pi blink
}
