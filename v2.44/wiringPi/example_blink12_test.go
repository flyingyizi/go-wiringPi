package wiringPi_test

import (
	"fmt"

	"github.com/flyingyizi/go-wiringPi/v2.44/wiringPi"
)

//wiringPi pin 0 is BCM_GPIO 17.
// Simple sequencer data
//	Triplets of LED, On/Off and delay
//Simple sequence over the first 12 GPIO pins - LEDs

var data = []uint{
	0, 1, 1,
	1, 1, 1,
	0, 0, 0, 2, 1, 1,
	1, 0, 0, 3, 1, 1,
	2, 0, 0, 4, 1, 1,
	3, 0, 0, 5, 1, 1,
	4, 0, 0, 6, 1, 1,
	5, 0, 0, 7, 1, 1,
	6, 0, 0, 11, 1, 1,
	7, 0, 0, 10, 1, 1,
	11, 0, 0, 13, 1, 1,
	10, 0, 0, 12, 1, 1,
	13, 0, 1,
	12, 0, 1,

	0, 0, 1, // Extra delay

	// Back again

	12, 1, 1,
	13, 1, 1,
	12, 0, 0, 10, 1, 1,
	13, 0, 0, 11, 1, 1,
	10, 0, 0, 7, 1, 1,
	11, 0, 0, 6, 1, 1,
	7, 0, 0, 5, 1, 1,
	6, 0, 0, 4, 1, 1,
	5, 0, 0, 3, 1, 1,
	4, 0, 0, 2, 1, 1,
	3, 0, 0, 1, 1, 1,
	2, 0, 0, 0, 1, 1,
	1, 0, 1,
	0, 0, 1,

	0, 0, 1, // Extra delay

	0, 9, 0, // End marker

}

func Example_blink12() {

	fmt.Println("Raspberry Pi - 12-LED Sequence")
	fmt.Println("==============================")
	fmt.Println("Connect LEDs up to the first 8 GPIO pins, then pins 11, 10, 13, 12 in")
	fmt.Println("    that order, then sit back and watch the show!")
	wiringPi.Setup()
	for i := 0; i < 14; i++ {
		wiringPi.SetPinMode(i, wiringPi.OutputPinMode)
	}

	var dataPtr, l, s, d uint

	for x := 0; x < 10; x++ {
		l = data[dataPtr] // LED
		dataPtr++
		s = data[dataPtr] // State
		dataPtr++
		d = data[dataPtr] // Duration (10ths)
		dataPtr++

		if s == 9 {
			dataPtr = 0
			continue
		} // 9 -> End Marker

		wiringPi.DigitalWrite(int(l), int(s))
		wiringPi.Delay(d * 100) //ms
	}
}
