package wiringPi_test

import "fmt"
import "github.com/flyingyizi/go-wiringPi/v2.44/wiringPi"

// LED Pin - wiringPi pin 0 is BCM_GPIO 17.

func Example_serial() {
	fd := wiringPi.SerialOpen("/dev/ttyAMA0", 115200)
	if fd < 0 {
		fmt.Println("Unable to open serial device: ")
		return
	}

	wiringPi.Setup()

	for i := 0; i < 10; i++ {
		wiringPi.DigitalWrite(LED, wiringPi.HIGH)
		wiringPi.Delay(500) //ms
		wiringPi.DigitalWrite(LED, wiringPi.LOW)
		wiringPi.Delay(500) //ms
	}

	// Output: Raspberry Pi blink
}

/*

int main ()
{
  int fd ;
  int count ;
  unsigned int nextTime ;

  if ((fd = serialOpen ("/dev/ttyAMA0", 115200)) < 0)
  {
    fprintf (stderr, "Unable to open serial device: %s\n", strerror (errno)) ;
    return 1 ;
  }

  if (wiringPiSetup () == -1)
  {
    fprintf (stdout, "Unable to start wiringPi: %s\n", strerror (errno)) ;
    return 1 ;
  }

  nextTime = millis () + 300 ;

  for (count = 0 ; count < 256 ; )
  {
    if (millis () > nextTime)
    {
      printf ("\nOut: %3d: ", count) ;
      fflush (stdout) ;
      serialPutchar (fd, count) ;
      nextTime += 300 ;
      ++count ;
    }

    delay (3) ;

    while (serialDataAvail (fd))
    {
      printf (" -> %3d", serialGetchar (fd)) ;
      fflush (stdout) ;
    }
  }

  printf ("\n") ;
  return 0 ;
}
*/
