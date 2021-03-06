//+build MEMMAP

package gpio

import (
	"errors"
	"os"
	"reflect"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/flyingyizi/go-wiringPi/board"
)

// Peripheral Offsets for the Raspberry Pi
const (
	ClkBase      = 0x00101000 // Clock registers
	GpioBase     = 0x00200000 // GPIO registers
	PwmBase      = 0x0020C000 // PWM registers
	PadsGpioBase = 0x00100000 //	PADS base
)

const SizeOfuint32 = 4 // bytes
const uint32BlockSize = SizeOfuint32 * 1024

var (
	gpioArry []uint32
	pwmArry  []uint32
	clkArry  []uint32
	padsArry []uint32

	memlock sync.Mutex

	gpio []byte
	clk  []byte
	pwm  []byte
	pads []byte
)

// Close unmaps GPIO memory
func gpioClose() (err error) {
	memlock.Lock()
	defer memlock.Unlock()

	err = syscall.Munmap(gpio)
	err = syscall.Munmap(pwm)
	err = syscall.Munmap(clk)
	err = syscall.Munmap(pads)
	return

}

func bytesToUint32Slince(b []byte) (data []uint32) {
	// Get the slice header
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&b))

	// The length and capacity of the slice are different.
	header.Len /= SizeOfuint32
	header.Cap /= SizeOfuint32

	// Convert slice header to an []uint32
	data = *(*[]uint32)(unsafe.Pointer(&header))
	return
}

func gpioOpen() (err error) {

	_, piGpioBase, err := board.GetBoardInfo()
	if err != nil {
		return
	}

	// Set the offsets into the memory interface.
	GPIO_PADS := piGpioBase + PadsGpioBase
	GPIO_CLOCK_BASE := piGpioBase + ClkBase
	GPIO_BASE := piGpioBase + GpioBase
	//GPIO_TIMER := piGpioBase + 0x0000B000
	GPIO_PWM := piGpioBase + PwmBase

	//	Try /dev/mem. If that fails, then
	//	try /dev/gpiomem. If that fails then game over.
	file, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, 0660)
	if err != nil {
		file, err = os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0660) //|os.O_CLOEXEC
		if err != nil {
			return errors.New("can not open /dev/mem or /dev/gpiomem, maybe try sudo")
		}
	}
	//fd can be closed after memory mapping
	defer file.Close()

	//	GPIO:
	gpio, err = syscall.Mmap(int(file.Fd()), GPIO_BASE, uint32BlockSize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return errors.New("mmap (GPIO) failed")
	}
	gpioArry = bytesToUint32Slince(gpio)

	//	PWM
	pwm, err = syscall.Mmap(int(file.Fd()), GPIO_PWM, uint32BlockSize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return errors.New("mmap (PWM) failed")
	}
	pwmArry = bytesToUint32Slince(pwm)

	//	Clock control (needed for PWM)
	clk, err = syscall.Mmap(int(file.Fd()), GPIO_CLOCK_BASE, uint32BlockSize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return errors.New("mmap (CLOCK) failed")
	}
	clkArry = bytesToUint32Slince(clk)

	//	The drive pads
	pads, err = syscall.Mmap(int(file.Fd()), GPIO_PADS, uint32BlockSize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return errors.New("mmap (PADS) failed")
	}
	padsArry = bytesToUint32Slince(pads)
	return
}

// Read the state(0:low, 1:high) of a pin
func gpioReadPin(bcmNumber uint8) (value uint, err error) {
	// Input level register offset (13 / 14 depending on bank)
	//In the datasheet on page 96, we seet that the GPLEVn register is
	//located 13 or 14 32-bit registers further than the gpio base register. GPLEV0 STORE 0~31,GPLEV1 STORE 32~53,
	levelReg := (bcmNumber)/32 + 13

	if (gpioArry[levelReg] & (1 << uint8(pin))) != 0 {
		return 1, nil
	}

	return 0, nil
}

// gpiopinMode sets the direction of a given pin (Input(0) or Output(1))
func gpiopinMode(bcmNumber uint8, direction Direction) {

	//In the datasheet at page 91 we find that the GPFSEL registers are organised per 10 pins.
	//So one 32-bit register contains the setup bits for 10 pins. *gpio.addr + ((g))/10 is
	// the register address that contains the GPFSEL bits of the pin "g"
	// Pin fsel register, 0 or 1 depending on bank
	fsel := (bcmNumber) / 10
	//There are three GPFSEL bits per pin (000: input, 001: output). The location
	//of these three bits inside the GPFSEL register is given by ((g)%10)*3
	shift := ((bcmNumber) % 10) * 3
	memlock.Lock()
	defer memlock.Unlock()

	if direction == inDirection {
		gpioArry[fsel] = gpioArry[fsel] &^ (7 << shift) //7:0b111 - pinmode is 3 bits
	} else {
		//This is also the reason that the comment says to "always use INP_GPIO(x) before using
		//OUT_GPIO(x)". This way you are sure that the other 2 bits are 0, and justifies the
		//use of a OR operation here. If you don't do that, you are not sure those bits will
		//be zero and you might have given the pin "g" a different setup.
		gpioArry[fsel] = gpioArry[fsel] &^ (7 << shift)
		gpioArry[fsel] = (gpioArry[fsel] &^ (7 << shift)) | (1 << shift)
	}
	pin.direction = direction

	//#define INP_GPIO(g)   *(gpio.addr + ((g)/10)) &= ~(7<<(((g)%10)*3))
	//#define OUT_GPIO(g)   *(gpio.addr + ((g)/10)) |=  (1<<(((g)%10)*3))
}

// gpioWritePin sets a given pin High(1) or Low(0)
// by setting the clear or set registers respectively
func gpioWritePin(bcmNumber uint8, state uint) error {

	p := (bcmNumber)

	// Clear register, 10 / 11 depending on bank
	// Set register, 7 / 8 depending on bank
	//In the datasheet on page 90, we seet that the GPSET register is
	//located 10 32-bit registers further than the gpio base register. GPCLR0 STORE 0~31,GPCLR1 STORE 32~53,
	clearReg := p/32 + 10
	//In the datasheet on page 90, we seet that the GPSET register is
	//located 7 32-bit registers further than the gpio base register. GPSET0 STORE 0~31,GPSET1 STORE 32~53,
	setReg := p/32 + 7

	memlock.Lock()
	defer memlock.Unlock()

	if state == 0 {
		gpioArry[clearReg] = 1 << (p & 31)
	} else {
		gpioArry[setReg] = 1 << (p & 31)
	}
	return nil
}

func gpioPullMode(bcmNumber uint8, pull Pull) {
	// Pull up/down/off register has offset 38 / 39, pull is 37
	pullClkReg := (bcmNumber)/32 + 38
	pullReg := 37
	shift := ((bcmNumber) % 32) // get 0 or 1 bank

	memlock.Lock()
	defer memlock.Unlock()

	switch pull {
	case PullDown, PullUp:
		gpioArry[pullReg] = gpioArry[pullReg]&^3 | uint32(pull)
	case PullOff:
		gpioArry[pullReg] = gpioArry[pullReg] &^ 3
	}

	// Wait for value to clock in, this is ugly, sorry :(
	time.Sleep(time.Microsecond)

	gpioArry[pullClkReg] = 1 << shift

	// Wait for value to clock in
	time.Sleep(time.Microsecond)

	gpioArry[pullReg] = gpioArry[pullReg] &^ 3
	gpioArry[pullClkReg] = 0

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
