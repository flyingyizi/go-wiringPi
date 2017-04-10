//+build linux

package rpi

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

//http://xillybus.com/tutorials/device-tree-zynq-1
//https://github.com/stianeikeland/go-rpio/blob/master/rpio.go

// Raspberry Pi Revision :: Model
const RPI_MODEL_A uint = 0       //   "Model A",	//  0
const RPI_MODEL_B uint = 1       //   "Model B",	//  1
const RPI_MODEL_A_PLUS uint = 2  //   "Model A+",	//  2
const RPI_MODEL_B_PLUS uint = 3  //   "Model B+",	//  3
const RPI_MODEL_2B uint = 4      //   "Pi 2",	//  4
const RPI_MODEL_ALPHA uint = 5   //   "Alpha",	//  5
const RPI_MODEL_CM uint = 6      //   "CM",		//  6
const RPI_MODEL_UNKNOWN uint = 7 //   "Unknown07",	// 07
const RPI_MODEL_3B uint = 8      //   "Pi 3",	// 08
const RPI_MODEL_ZERO uint = 9    //   "Pi Zero",	// 09
const RPI_MODEL_CM3 uint = 10    //   "CM3",	// 10
const RPI_MODEL_ZERO_W uint = 12 //   "Pi Zero-W",	// 12

var RaspberryModel = map[uint]string{
	RPI_MODEL_A:       "Model A",   //  0
	RPI_MODEL_B:       "Model B",   //  1
	RPI_MODEL_A_PLUS:  "Model A+",  //  2
	RPI_MODEL_B_PLUS:  "Model B+",  //  3
	RPI_MODEL_2B:      "Pi 2",      //
	RPI_MODEL_ALPHA:   "Alpha",     //  5
	RPI_MODEL_CM:      "CM",        //  6
	RPI_MODEL_UNKNOWN: "Unknown07", // 07
	RPI_MODEL_3B:      "Pi 3",      // 08
	RPI_MODEL_ZERO:    "Pi Zero",   // 09
	RPI_MODEL_CM3:     "CM3",       // 10
	RPI_MODEL_ZERO_W:  "Pi Zero-W", // 12
}

const RPI_VERSION_1 uint = 0
const RPI_VERSION_1_1 uint = 1
const RPI_VERSION_1_2 uint = 2
const RPI_VERSION_2 uint = 3

const RPI_MAKER_SONY uint = 0
const RPI_MAKER_EGOMAN uint = 1
const RPI_MAKER_EMBEST uint = 2
const RPI_MAKER_UNKNOWN uint = 3

const SIZEOF_UINT32 = 4 // bytes
const uint32BlockSize = SIZEOF_UINT32 * 1024

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

type Pull uint8

// Pull Up / Down / Off
const (
	PullOff Pull = iota
	PullDown
	PullUp
)

type Pin uint8

// Set pin as Input
func (pin Pin) Input() {
	gpiopinMode(pin, 0)
}

// Set pin as Output
func (pin Pin) Output() {
	gpiopinMode(pin, 1)
}

// Set pin High
func (pin Pin) High() {
	gpioWritePin(pin, 1)
}

// Set pin Low
func (pin Pin) Low() {
	gpioWritePin(pin, 0)
}

// Close unmaps GPIO memory
func Close() (err error) {
	memlock.Lock()
	defer memlock.Unlock()

	err = syscall.Munmap(gpio)
	err = syscall.Munmap(pwm)
	err = syscall.Munmap(clk)
	err = syscall.Munmap(pads)
	return

}

func BytesToUint32Slince(b []byte) (data []uint32) {
	// Get the slice header
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&b))

	// The length and capacity of the slice are different.
	header.Len /= SIZEOF_UINT32
	header.Cap /= SIZEOF_UINT32

	// Convert slice header to an []uint32
	data = *(*[]uint32)(unsafe.Pointer(&header))
	return
}

func Init() (err error) {

	//fd can be closed after memory mapping
	defer file.Close()

	_, bmodel, _, _, _, _, err := PiBoardId()
	fmt.Println("modes is %s", RaspberryModel[bmodel])
	var piGpioBase int64 = 0x20000000
	if bmodel == RPI_MODEL_A || bmodel == RPI_MODEL_B || bmodel == RPI_MODEL_A_PLUS || bmodel == RPI_MODEL_B_PLUS || bmodel == RPI_MODEL_ALPHA || bmodel == RPI_MODEL_CM || bmodel == RPI_MODEL_ZERO || bmodel == RPI_MODEL_ZERO_W {
		// piGpioBase:
		//	The base address of the GPIO memory mapped hardware IO
		piGpioBase = 0x20000000

	} else {
		piGpioBase = 0x3F000000

	}
	// Set the offsets into the memory interface.
	GPIO_PADS := piGpioBase + 0x00100000
	GPIO_CLOCK_BASE := piGpioBase + 0x00101000
	GPIO_BASE := piGpioBase + 0x00200000
	//GPIO_TIMER := piGpioBase + 0x0000B000
	GPIO_PWM := piGpioBase + 0x0020C000

	//	Try /dev/mem. If that fails, then
	//	try /dev/gpiomem. If that fails then game over.
	file, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, 0660)
	if err != nil {
		file, err = os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0660) //|os.O_CLOEXEC
		if err != nil {
			return errors.New("can not open /dev/mem or /dev/gpiomem, maybe try sudo")
		}
	}

	//	GPIO:
	gpio, err = syscall.Mmap(int(file.Fd()), GPIO_BASE, uint32BlockSize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return errors.New("mmap (GPIO) failed")
	}
	gpioArry = BytesToUint32Slince(gpio)

	//	PWM
	pwm, err = syscall.Mmap(int(file.Fd()), GPIO_PWM, uint32BlockSize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return errors.New("mmap (PWM) failed")
	}
	pwmArry = BytesToUint32Slince(pwm)

	//	Clock control (needed for PWM)
	clk, err = syscall.Mmap(int(file.Fd()), GPIO_CLOCK_BASE, uint32BlockSize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return errors.New("mmap (CLOCK) failed")
	}
	clkArry = BytesToUint32Slince(clk)

	//	The drive pads
	pads, err = syscall.Mmap(int(file.Fd()), GPIO_PADS, uint32BlockSize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return errors.New("mmap (PADS) failed")
	}
	padsArry = BytesToUint32Slince(pads)
	return
}

// Read the state(0:low, 1:high) of a pin
func gpioReadPin(pin Pin) int {
	// Input level register offset (13 / 14 depending on bank)
	//In the datasheet on page 96, we seet that the GPLEVn register is
	//located 13 or 14 32-bit registers further than the gpio base register. GPLEV0 STORE 0~31,GPLEV1 STORE 32~53,

	levelReg := uint8(pin)/32 + 13

	if (gpioArry[levelReg] & (1 << uint8(pin))) != 0 {
		return 1
	}

	return 0
}

// Toggle a pin state (high -> low -> high)
func (pin Pin) TogglePin() {
	switch gpioReadPin(pin) {
	case 0:
		pin.High()
	default:
		pin.Low()
	}
}

// gpiopinMode sets the direction of a given pin (Input(0) or Output(1))
func gpiopinMode(pin Pin, direction int) {

	//In the datasheet at page 91 we find that the GPFSEL registers are organised per 10 pins.
	//So one 32-bit register contains the setup bits for 10 pins. *gpio.addr + ((g))/10 is
	// the register address that contains the GPFSEL bits of the pin "g"
	// Pin fsel register, 0 or 1 depending on bank
	fsel := uint8(pin) / 10
	//There are three GPFSEL bits per pin (000: input, 001: output). The location
	//of these three bits inside the GPFSEL register is given by ((g)%10)*3
	shift := (uint8(pin) % 10) * 3
	memlock.Lock()
	defer memlock.Unlock()

	if direction == 0 {
		gpioArry[fsel] = gpioArry[fsel] &^ (7 << shift) //7:0b111 - pinmode is 3 bits
	} else {
		//This is also the reason that the comment says to "always use INP_GPIO(x) before using
		//OUT_GPIO(x)". This way you are sure that the other 2 bits are 0, and justifies the
		//use of a OR operation here. If you don't do that, you are not sure those bits will
		//be zero and you might have given the pin "g" a different setup.
		gpioArry[fsel] = gpioArry[fsel] &^ (7 << shift)
		gpioArry[fsel] = (gpioArry[fsel] &^ (7 << shift)) | (1 << shift)
	}

	//#define INP_GPIO(g)   *(gpio.addr + ((g)/10)) &= ~(7<<(((g)%10)*3))
	//#define OUT_GPIO(g)   *(gpio.addr + ((g)/10)) |=  (1<<(((g)%10)*3))
}

// gpioWritePin sets a given pin High(1) or Low(0)
// by setting the clear or set registers respectively
func gpioWritePin(pin Pin, state int) {

	p := uint8(pin)

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

}

func gpioPullMode(pin Pin, pull Pull) {
	// Pull up/down/off register has offset 38 / 39, pull is 37
	pullClkReg := uint8(pin)/32 + 38
	pullReg := 37
	shift := (uint8(pin) % 32) // get 0 or 1 bank

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

func piGPIOLayout() (err error) {
	cpuinfo, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return err

	}
	lines := strings.Split(string(cpuinfo), "\n")

	str := `Unable to determine hardware version. I see: %s 
     - expecting BCM2708, BCM2709 or BCM2835. 
    If this is a genuine Raspberry Pi then please report this 
    to projects@drogon.net. If this is not a Raspberry Pi then you 
    are on your own as wiringPi is designed to support the 
    Raspberry Pi ONLY.\n`
	var ErrHardWare error = errors.New(str)

	for _, line := range lines {
		fields := strings.Split(line, ":")
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		if key == "Hardware" {
			if value == "BCM2708" || value == "BCM2709" || value == "BCM2835" {
				ErrHardWare = nil
			}
		} else if key == "Revision" {

			return ErrHardWare
		}
		//unicode.IsNumber

	}
	return ErrHardWare
}

func PiBoardId() (pcbrev uint, bmodel uint, processor uint, manufacturer uint, ram uint, bWarranty uint, err error) {

	str := `Unable to determine boardinfo. If this is not a Raspberry Pi then you 
    are on your own as wiringPi is designed to support the 
    Raspberry Pi ONLY.\n`
	var ErrRevision = errors.New(str)

	cpuinfo, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		err = ErrRevision
		return

	}
	lines := strings.Split(string(cpuinfo), "\n")

	revisionValue := ""
	for _, l := range lines {
		fields := strings.Split(l, ":")
		if len(fields) == 2 {
			key := strings.TrimSpace(fields[0])
			value := strings.TrimSpace(fields[1])
			if key == "Revision" {
				ErrRevision = nil
				revisionValue = value
				break
			}
		}
	}
	if ErrRevision != nil {
		return 0, 0, 0, 0, 0, 0, ErrRevision
	}

	// If longer than 4, we'll assume it's been overvolted
	if len(revisionValue) > 4 {
		bWarranty = 1
		// Extract last 4 characters
		revisionValue = revisionValue[len(revisionValue)-4:]
	}

	// Hex number with no leading 0x
	i, err := strconv.ParseUint(revisionValue, 16, 32)
	revision := (uint)(i)
	if err != nil {
		return
	}

	// SEE: https://github.com/AndrewFromMelbourne/raspberry_pi_revision
	scheme := (revision & (1 << 23)) >> 23

	if scheme > 0 {
		pcbrev = (revision & (0x0F << 0)) >> 0
		bmodel = (revision & (0xFF << 4)) >> 4
		processor = (revision & (0x0F << 12)) >> 12 // Not used for now.
		manufacturer = (revision & (0x0F << 16)) >> 16
		ram = (revision & (0x07 << 20)) >> 20
		bWarranty = (revision & (0x03 << 24)) >> 24

	} else {
		switch revisionValue {
		case "0002":
			bmodel = RPI_MODEL_B
			pcbrev = RPI_VERSION_1
			ram = 0
			manufacturer = RPI_MAKER_EGOMAN
		case "0003":
			bmodel = RPI_MODEL_B
			pcbrev = RPI_VERSION_1_1
			ram = 0
			manufacturer = RPI_MAKER_EGOMAN
		case "0004":
			bmodel = RPI_MODEL_B
			pcbrev = RPI_VERSION_1_2
			ram = 0
			manufacturer = RPI_MAKER_SONY
		case "0005":
			fallthrough
		case "0006":
			fallthrough
		case "000f":
			fallthrough
		case "000d":
			bmodel = RPI_MODEL_B
			pcbrev = RPI_VERSION_1_2
			ram = 0
			manufacturer = RPI_MAKER_EGOMAN
		case "0007":
			fallthrough
		case "0009":
			bmodel = RPI_MODEL_A
			pcbrev = RPI_VERSION_1_2
			ram = 0
			manufacturer = RPI_MAKER_EGOMAN
		case "0008":
			bmodel = RPI_MODEL_A
			pcbrev = RPI_VERSION_1_2
			ram = 0
			manufacturer = RPI_MAKER_SONY
		case "0010":
			fallthrough
		case "0016":
			bmodel = RPI_MODEL_B_PLUS
			pcbrev = RPI_VERSION_1_2
			ram = 1
			manufacturer = RPI_MAKER_SONY
		case "0013":
			bmodel = RPI_MODEL_B_PLUS
			pcbrev = RPI_VERSION_1_2
			ram = 1
			manufacturer = RPI_MAKER_EMBEST
		case "0019":
			bmodel = RPI_MODEL_B_PLUS
			pcbrev = RPI_VERSION_1_2
			ram = 1
			manufacturer = RPI_MAKER_EGOMAN
		case "0011":
			fallthrough
		case "0017":
			bmodel = RPI_MODEL_CM
			pcbrev = RPI_VERSION_1_1
			ram = 1
			manufacturer = RPI_MAKER_SONY
		case "0014":
			bmodel = RPI_MODEL_CM
			pcbrev = RPI_VERSION_1_1
			ram = 1
			manufacturer = RPI_MAKER_EMBEST
		case "001a":
			bmodel = RPI_MODEL_CM
			pcbrev = RPI_VERSION_1_1
			ram = 1
			manufacturer = RPI_MAKER_EGOMAN
		case "0012":
			fallthrough
		case "0018":
			bmodel = RPI_MODEL_A_PLUS
			pcbrev = RPI_VERSION_1_1
			ram = 0
			manufacturer = RPI_MAKER_SONY
		case "0015":
			bmodel = RPI_MODEL_A_PLUS
			pcbrev = RPI_VERSION_1_1
			ram = 1
			manufacturer = RPI_MAKER_EMBEST
		case "001b":
			bmodel = RPI_MODEL_A_PLUS
			pcbrev = RPI_VERSION_1_1
			ram = 0
			manufacturer = RPI_MAKER_EGOMAN

		}

	}

	return

	//-------------------------------------------------------------------------
	// SEE: https://github.com/AndrewFromMelbourne/raspberry_pi_revision
	//-------------------------------------------------------------------------
	//
	// The file /proc/cpuinfo contains a line such as:-
	//
	// Revision    : 0003
	//
	// that holds the revision number of the Raspberry Pi.
	// Known revisions (prior to the Raspberry Pi 2) are:
	//
	//     +----------+---------+---------+--------+-------------+
	//     | Revision |  Model  | PCB Rev | Memory | Manufacture |
	//     +----------+---------+---------+--------+-------------+
	//     |   0000   |         |         |        |             |
	//     |   0001   |         |         |        |             |
	//     |   0002   |    B    |    1    | 256 MB |             |
	//     |   0003   |    B    |    1    | 256 MB |             |
	//     |   0004   |    B    |    2    | 256 MB |   Sony      |
	//     |   0005   |    B    |    2    | 256 MB |   Qisda     |
	//     |   0006   |    B    |    2    | 256 MB |   Egoman    |
	//     |   0007   |    A    |    2    | 256 MB |   Egoman    |
	//     |   0008   |    A    |    2    | 256 MB |   Sony      |
	//     |   0009   |    A    |    2    | 256 MB |   Qisda     |
	//     |   000a   |         |         |        |             |
	//     |   000b   |         |         |        |             |
	//     |   000c   |         |         |        |             |
	//     |   000d   |    B    |    2    | 512 MB |   Egoman    |
	//     |   000e   |    B    |    2    | 512 MB |   Sony      |
	//     |   000f   |    B    |    2    | 512 MB |   Qisda     |
	//     |   0010   |    B+   |    1    | 512 MB |   Sony      |
	//     |   0011   | compute |    1    | 512 MB |   Sony      |
	//     |   0012   |    A+   |    1    | 256 MB |   Sony      |
	//     |   0013   |    B+   |    1    | 512 MB |   Embest    |
	//     |   0014   | compute |    1    | 512 MB |   Sony      |
	//     |   0015   |    A+   |    1    | 256 MB |   Sony      |
	//     +----------+---------+---------+--------+-------------+
	//
	// If the Raspberry Pi has been over-volted (voiding the warranty) the
	// revision number will have 100 at the front. e.g. 1000002.
	//
	//-------------------------------------------------------------------------
	//
	// With the release of the Raspberry Pi 2, there is a new encoding of the
	// Revision field in /proc/cpuinfo. The bit fields are as follows
	//
	//     +----+----+----+----+----+----+----+----+
	//     |FEDC|BA98|7654|3210|FEDC|BA98|7654|3210|
	//     +----+----+----+----+----+----+----+----+
	//     |    |    |    |    |    |    |    |AAAA|
	//     |    |    |    |    |    |BBBB|BBBB|    |
	//     |    |    |    |    |CCCC|    |    |    |
	//     |    |    |    |DDDD|    |    |    |    |
	//     |    |    | EEE|    |    |    |    |    |
	//     |    |    |F   |    |    |    |    |    |
	//     |    |   G|    |    |    |    |    |    |
	//     |    |  H |    |    |    |    |    |    |
	//     +----+----+----+----+----+----+----+----+
	//     |1098|7654|3210|9876|5432|1098|7654|3210|
	//     +----+----+----+----+----+----+----+----+
	//
	// +---+-------+--------------+--------------------------------------------+
	// | # | bits  |   contains   | values                                     |
	// +---+-------+--------------+--------------------------------------------+
	// | A | 00-03 | PCB Revision | (the pcb revision number)                  |
	// | B | 04-11 | Model name   | A, B, A+, B+, B Pi2, Alpha, Compute Module |
	// |   |       |              | unknown, B Pi3, Zero                       |
	// | C | 12-15 | Processor    | BCM2835, BCM2836, BCM2837                  |
	// | D | 16-19 | Manufacturer | Sony, Egoman, Embest, unknown, Embest      |
	// | E | 20-22 | Memory size  | 256 MB, 512 MB, 1024 MB                    |
	// | F | 23-23 | encoded flag | (if set, revision is a bit field)          |
	// | G | 24-24 | waranty bit  | (if set, warranty void - Pre Pi2)          |
	// | H | 25-25 | waranty bit  | (if set, warranty void - Post Pi2)         |
	// +---+-------+--------------+--------------------------------------------+
	//
	// Also, due to some early issues the warranty bit has been move from bit
	// 24 to bit 25 of the revision number (i.e. 0x2000000).

}
