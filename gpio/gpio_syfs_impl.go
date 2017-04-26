package gpio

import (
	"fmt"
	"os"
	"strconv"
)

var fs = [64]*os.File{nil}

//ref https://github.com/brian-armstrong/gpio/blob/master/sysfs.go

// Read the state(0:low, 1:high) of a pin
func gpioReadPin(bcmNumber uint8) (value uint, err error) {
	file := fs[bcmNumber]
	file.Seek(0, 0)
	buf := make([]byte, 1)
	_, err = file.Read(buf)
	if err != nil {
		return 0, err
	}
	c := buf[0]
	switch c {
	case '0':
		return 0, nil
	case '1':
		return 1, nil
	default:
		return 0, fmt.Errorf("read inconsistent value in pinfile, %c", c)
	}
}

// gpioWritePin sets a given pin High(1) or Low(0)
// by setting the clear or set registers respectively
func gpioWritePin(bcmNumber uint8, state int) error {
	var buf []byte
	switch state {
	case 0:
		buf = []byte{'0'}
	case 1:
		buf = []byte{'1'}
	default:
		return fmt.Errorf("invalid output value %d", state)
	}
	file := fs[bcmNumber]

	_, err := file.Write(buf)
	return err
}

func openPin(bcmNumber uint8, write bool) error {
	flags := os.O_RDONLY
	if write {
		flags = os.O_RDWR
	}
	f, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", bcmNumber), flags, 0600)
	if err != nil {
		return fmt.Errorf("failed to open gpio %d value file for reading\n", bcmNumber)
	}
	fs[bcmNumber] = f
	return nil
}

// Close unmaps GPIO memory
func gpioClose() (err error) {
	//todo
	return

}

func gpioOpen() (err error) {

	//todo
	return
}

// gpiopinMode sets the direction of a given pin (Input(0) or Output(1))
func gpiopinMode(bcmNumber uint8, direction Direction) {

	/*
	   # Set up GPIO 4 and set to output
	   echo "4" > /sys/class/gpio/export
	   echo "out" > /sys/class/gpio/gpio4/direction

	   # Set up GPIO 7 and set to input
	   echo "7" > /sys/class/gpio/export
	   echo "in" > /sys/class/gpio/gpio7/direction

	   # Write output
	   echo "1" > /sys/class/gpio/gpio4/value

	   # Read from input
	   cat /sys/class/gpio/gpio7/value

	   # Clean up
	   echo "4" > /sys/class/gpio/unexport
	   echo "7" > /sys/class/gpio/unexport


	*/

	dir, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", bcmNumber), os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("failed to open gpio %d direction file for writing\n", bcmNumber)
		os.Exit(1)
	}
	defer dir.Close()

	if direction == inDirection {
		dir.Write([]byte("in"))
	} else {
		dir.Write([]byte("out"))
	}
}

func exportGPIO(p Pin) {
	export, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("failed to open gpio export file for writing\n")
		os.Exit(1)
	}
	defer export.Close()
	export.Write([]byte(strconv.Itoa(int(p.Number))))
}

func unexportGPIO(p Pin) {
	export, err := os.OpenFile("/sys/class/gpio/unexport", os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("failed to open gpio unexport file for writing\n")
		os.Exit(1)
	}
	defer export.Close()
	export.Write([]byte(strconv.Itoa(int(p.Number))))
}

func gpioPullMode(bcmNumber uint8, pull Pull) {
	// Pull up/down/off register has offset 38 / 39, pull is 37
	//todo
}

/*
https://github.com/jameswalmsley/RaspberryPi-FreeRTOS/blob/master/Demo/Drivers/gpio.c
typedef struct {
	unsigned long	GPFSEL[6];	///< Function selection registers.
	unsigned long	Reserved_1;
	unsigned long	GPSET[2];
	unsigned long	Reserved_2;
	unsigned long	GPCLR[2];
	unsigned long	Reserved_3;
	unsigned long	GPLEV[2];
	unsigned long	Reserved_4;
	unsigned long	GPEDS[2];
	unsigned long	Reserved_5;
	unsigned long	GPREN[2];
	unsigned long	Reserved_6;
	unsigned long	GPFEN[2];
	unsigned long	Reserved_7;
	unsigned long	GPHEN[2];
	unsigned long	Reserved_8;
	unsigned long	GPLEN[2];
	unsigned long	Reserved_9;
	unsigned long	GPAREN[2];
	unsigned long	Reserved_A;
	unsigned long	GPAFEN[2];
	unsigned long	Reserved_B;
	unsigned long	GPPUD[1];
	unsigned long	GPPUDCLK[2];
	//Ignoring the reserved and test bytes
} BCM2835_GPIO_REGS;

*/
