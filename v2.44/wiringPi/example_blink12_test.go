package wiringPi_test

import (
	"fmt"

	"github.com/flyingyizi/go-wiringPi/v2.44/wiringPi"
)

//wiringPi pin 0 is BCM_GPIO 17.
// Simple sequencer data
//	Triplets of LED, On/Off and delay
//Simple sequence over the first 12 GPIO pins - LEDs

var data []int = {
            0, 1, 1,
            1, 1, 1,
  0, 0, 0,  2, 1, 1,
  1, 0, 0,  3, 1, 1,
  2, 0, 0,  4, 1, 1,
  3, 0, 0,  5, 1, 1,
  4, 0, 0,  6, 1, 1,
  5, 0, 0,  7, 1, 1,
  6, 0, 0, 11, 1, 1,
  7, 0, 0, 10, 1, 1,
 11, 0, 0, 13, 1, 1,
 10, 0, 0, 12, 1, 1,
 13, 0, 1,
 12, 0, 1,

  0, 0, 1,	// Extra delay

// Back again

           12, 1, 1,
           13, 1, 1,
 12, 0, 0, 10, 1, 1,
 13, 0, 0, 11, 1, 1,
 10, 0, 0,  7, 1, 1,
 11, 0, 0,  6, 1, 1,
  7, 0, 0,  5, 1, 1,
  6, 0, 0,  4, 1, 1,
  5, 0, 0,  3, 1, 1,
  4, 0, 0,  2, 1, 1,
  3, 0, 0,  1, 1, 1,
  2, 0, 0,  0, 1, 1,
  1, 0, 1,
  0, 0, 1,

  0, 0, 1,	// Extra delay

  0, 9, 0,	// End marker

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

	/*
	  for (;;)
	  {
	    l = data [dataPtr++] ;	// LED
	    s = data [dataPtr++] ;	// State
	    d = data [dataPtr++] ;	// Duration (10ths)

	    if (s == 9)			// 9 -> End Marker
	    {
	      dataPtr = 0 ;
	      continue ;
	    }

	    digitalWrite (l, s) ;
	    delay        (d * 100) ;
	  }

	*/
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

/*

int data [] =
{
            0, 1, 1,
            1, 1, 1,
  0, 0, 0,  2, 1, 1,
  1, 0, 0,  3, 1, 1,
  2, 0, 0,  4, 1, 1,
  3, 0, 0,  5, 1, 1,
  4, 0, 0,  6, 1, 1,
  5, 0, 0,  7, 1, 1,
  6, 0, 0, 11, 1, 1,
  7, 0, 0, 10, 1, 1,
 11, 0, 0, 13, 1, 1,
 10, 0, 0, 12, 1, 1,
 13, 0, 1,
 12, 0, 1,

  0, 0, 1,	// Extra delay

// Back again

           12, 1, 1,
           13, 1, 1,
 12, 0, 0, 10, 1, 1,
 13, 0, 0, 11, 1, 1,
 10, 0, 0,  7, 1, 1,
 11, 0, 0,  6, 1, 1,
  7, 0, 0,  5, 1, 1,
  6, 0, 0,  4, 1, 1,
  5, 0, 0,  3, 1, 1,
  4, 0, 0,  2, 1, 1,
  3, 0, 0,  1, 1, 1,
  2, 0, 0,  0, 1, 1,
  1, 0, 1,
  0, 0, 1,

  0, 0, 1,	// Extra delay

  0, 9, 0,	// End marker

} ;


int main (void)
{
  int pin ;
  int dataPtr ;
  int l, s, d ;

  printf ("Raspberry Pi - 12-LED Sequence\n") ;
  printf ("==============================\n") ;
  printf ("\n") ;
  printf ("Connect LEDs up to the first 8 GPIO pins, then pins 11, 10, 13, 12 in\n") ;
  printf ("    that order, then sit back and watch the show!\n") ;

  wiringPiSetup () ;

  for (pin = 0 ; pin < 14 ; ++pin)
    pinMode (pin, OUTPUT) ;

  dataPtr = 0 ;

  for (;;)
  {
    l = data [dataPtr++] ;	// LED
    s = data [dataPtr++] ;	// State
    d = data [dataPtr++] ;	// Duration (10ths)

    if (s == 9)			// 9 -> End Marker
    {
      dataPtr = 0 ;
      continue ;
    }

    digitalWrite (l, s) ;
    delay        (d * 100) ;
  }

  return 0 ;
}

*/
