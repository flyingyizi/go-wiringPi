
package wiringPi_test


import  a "github.com/flyingyizi/go-wiringPi/v2.44/wiringPi"

// LED Pin - wiringPi pin 0 is BCM_GPIO 17.
const LED int=0

func Example()  {
fmt.Println("Raspberry Pi blink\n")
a.
}




/*

#define	LED	0

int main (void)
{
  printf ("Raspberry Pi blink\n") ;

  wiringPiSetup () ;
  pinMode (LED, OUTPUT) ;

  for (;;)
  {
    digitalWrite (LED, HIGH) ;	// On
    delay (500) ;		// mS
    digitalWrite (LED, LOW) ;	// Off
    delay (500) ;
  }
  return 0 ;
}

*/
