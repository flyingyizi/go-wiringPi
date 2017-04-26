package gpio

import "errors"

type Pull uint8

// Pull Up / Down / Off
const (
	PullOff Pull = iota
	PullDown
	PullUp
)

type Direction uint

const (
	inDirection Direction = iota
	outDirection
)

// Pin represents a single pin, which can be used either for reading or writing
type Pin struct {
	bcmNumber uint8
	direction Direction
}

// Set pin as Input
func (pin Pin) Input() {
	gpiopinMode(pin.bcmNumber, inDirection)
}

// Set pin as Output
func (pin Pin) Output() {
	gpiopinMode(pin.bcmNumber, outDirection)
}

// High sets the value of an output pin to logic high
func (p Pin) High() error {
	if p.direction != outDirection {
		return errors.New("pin is not configured for output")
	}
	return gpioWritePin(p.bcmNumber, 1)
}

// Low sets the value of an output pin to logic low
func (p Pin) Low() error {
	if p.direction != outDirection {
		return errors.New("pin is not configured for output")
	}
	return gpioWritePin(p.bcmNumber, 0)
}

// Toggle a pin state (high -> low -> high)
func (pin Pin) TogglePin() {
	value, _ := gpioReadPin(pin.bcmNumber)
	switch value {
	case 0:
		pin.High()
	default:
		pin.Low()
	}
}

func (pin Pin) Read() (value uint, err error) {

	if pin.direction != inDirection {
		return 0, errors.New("pin is not configured for input")
	}
	value, err = gpioReadPin(pin.bcmNumber)
	return

}

// Close unmaps GPIO memory
func Close() (err error) {
	return gpioClose()

}

func Open() (err error) {

	return gpioOpen()
}

/*
const (
	edgeNone edge = iota
	edgeRising
	edgeFalling
	edgeBoth
)
func setEdgeTrigger(p Pin, e edge) {
	edge, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/edge", p.Number), os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("failed to open gpio %d edge file for writing\n", p.Number)
		os.Exit(1)
	}
	defer edge.Close()

	switch e {
	case edgeNone:
		edge.Write([]byte("none"))
	case edgeRising:
		edge.Write([]byte("rising"))
	case edgeFalling:
		edge.Write([]byte("falling"))
	case edgeBoth:
		edge.Write([]byte("both"))
	default:
		panic(fmt.Sprintf("setEdgeTrigger called with invalid edge %d", e))
	}
}

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
