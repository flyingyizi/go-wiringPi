package wiringPi_test

import (
	"fmt"

	"github.com/flyingyizi/go-wiringPi/v2.44/wiringPi"
)

//wiringPi pin 0 is BCM_GPIO 17.
//Simple sequence over the first 8 GPIO pins - LEDs

func Example_blink8() {
	fmt.Println("Raspberry Pi - 8-LED Sequencer")
	fmt.Println("==============================")
	fmt.Println("Connect LEDs to the first 8 GPIO pins and watch ...")

	wiringPi.Setup()
	for i := 0; i < 8; i++ {
		wiringPi.SetPinMode(i, wiringPi.OutputPinMode)
	}

	for i := 0; i < 10; i++ {
		for led := 0; led < 8; led++ {
			wiringPi.DigitalWrite(led, wiringPi.HIGH)
			wiringPi.Delay(100) //ms

		}
		for led := 0; led < 8; led++ {
			wiringPi.DigitalWrite(led, wiringPi.LOW)
			wiringPi.Delay(100) //ms

		}
	}
}
