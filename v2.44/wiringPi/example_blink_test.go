package wiringPi_test

import "fmt"
import "github.com/flyingyizi/go-wiringPi/v2.44/wiringPi"

// LED Pin - wiringPi pin 0 is BCM_GPIO 17.
const LED int = 0

func Example_blink() {
	fmt.Println("Raspberry Pi blink")
	wiringPi.Setup()
	wiringPi.SetPinMode(LED, wiringPi.OutputPinMode)

	for i := 0; i < 10; i++ {
		wiringPi.DigitalWrite(LED, wiringPi.HIGH)
		wiringPi.Delay(500) //ms
		wiringPi.DigitalWrite(LED, wiringPi.LOW)
		wiringPi.Delay(500) //ms
	}

	//// Output: Raspberry Pi blink
}
