package main

_ import "github.com/flyingyizi/go-wiringPi/spi"


const (
SPI_CHAN		=0
NUM_TIMES		=100
MAX_SIZE		=(1024*1024)
)	

/*
func main() {
	d, err := i2c.Open(0x39)
	if err != nil {
		panic(err)
	}

	// opens a 10-bit address
	d, err = i2c.Open(i2c.TenBit(0x78))

	if err != nil {
		panic(err)
	}

	d.Close()

	_ = d
}





void spiSetup (int speed)
{
  if ((myFd = wiringPiSPISetup (SPI_CHAN, speed)) < 0)
  {
    fprintf (stderr, "Can't open the SPI bus: %s\n", strerror (errno)) ;
    exit (EXIT_FAILURE) ;
  }
}


int main (void)
{
  int speed, times, size ;
  unsigned int start, end ;
  int spiFail ;
  unsigned char *myData ;
  double timePerTransaction, perfectTimePerTransaction, dataSpeed ;

  if ((myData = malloc (MAX_SIZE)) == NULL)
  {
    fprintf (stderr, "Unable to allocate buffer: %s\n", strerror (errno)) ;
    exit (EXIT_FAILURE) ;
  }

  wiringPiSetup () ;

  for (speed = 1 ; speed <= 32 ; speed *= 2)
  {
    printf ("+-------+--------+----------+----------+-----------+------------+\n") ;
    printf ("|   MHz |   Size | mS/Trans |      TpS |    Mb/Sec | Latency mS |\n") ;
    printf ("+-------+--------+----------+----------+-----------+------------+\n") ;

    spiFail = FALSE ;
    spiSetup (speed * 1000000) ;
    for (size = 1 ; size <= MAX_SIZE ; size *= 2)
    {
      printf ("| %5d | %6d ", speed, size) ;

      start = millis () ;
      for (times = 0 ; times < NUM_TIMES ; ++times)
	if (wiringPiSPIDataRW (SPI_CHAN, myData, size) == -1)
	{
	  printf ("SPI failure: %s\n", strerror (errno)) ;
	  spiFail = TRUE ;
	  break ;
	}
      end = millis () ;

      if (spiFail)
	break ;

      timePerTransaction        = ((double)(end - start) / (double)NUM_TIMES) / 1000.0 ;
      dataSpeed                 =  (double)(size * 8)    / (1024.0 * 1024.0) / timePerTransaction  ;
      perfectTimePerTransaction = ((double)(size * 8))   / ((double)(speed * 1000000)) ;

      printf ("| %8.3f ", timePerTransaction * 1000.0) ;
      printf ("| %8.1f ", 1.0 / timePerTransaction) ;
      printf ("| %9.5f ", dataSpeed) ;
      printf ("|   %8.5f ", (timePerTransaction - perfectTimePerTransaction) * 1000.0) ;
      printf ("|\n") ;

    }

    close (myFd) ;
    printf ("+-------+--------+----------+----------+-----------+------------+\n") ;
    printf ("\n") ;
  }

  return 0 ;
}
*/