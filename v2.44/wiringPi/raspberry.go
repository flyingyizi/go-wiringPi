//+build linux

package wiringPi

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

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

const RPI_VERSION_1 uint = 0
const RPI_VERSION_1_1 uint = 1
const RPI_VERSION_1_2 uint = 2
const RPI_VERSION_2 uint = 3

const RPI_MAKER_SONY uint = 0
const RPI_MAKER_EGOMAN uint = 1
const RPI_MAKER_EMBEST uint = 2
const RPI_MAKER_UNKNOWN uint = 3

const uint32BlockSize = 4 * 1024

var (
	gpioArry []uint32
	pwmArry  []uint32
	clkArr   []uint32
	padsArry []uint32
)

func Init() (err error) {
	// piGpioBase:
	//	The base address of the GPIO memory mapped hardware IO
	piGpioBase := int64(0x20000000)

	//	Try /dev/mem. If that fails, then
	//	try /dev/gpiomem. If that fails then game over.
	file, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, 0660)
	if err != nil {
		file, err = os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0660) //|os.O_CLOEXEC

		return errors.New("can not open /dev/mem  or /dev/gpiomem, maybe try sudo")
	}
	//fd can be closed after memory mapping
	defer file.Close()

	_, bmodel, _, _, _, _, err := piBoardId()

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

	//	GPIO:
	gpio, err := syscall.Mmap(int(file.Fd()), GPIO_BASE, uint32BlockSize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return errors.New("mmap (GPIO) failed")
	}
	i := (*[uint32BlockSize / 4]uint32)(unsafe.Pointer(&gpio[0]))
	gpioArry = i[:]
	//	PWM
	pwm, err := syscall.Mmap(int(file.Fd()), GPIO_PWM, uint32BlockSize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return errors.New("mmap (PWM) failed")
	}
	i = (*[uint32BlockSize / 4]uint32)(unsafe.Pointer(&pwm[0]))
	pwmArry = i[:]

	//	Clock control (needed for PWM)
	clk, err := syscall.Mmap(int(file.Fd()), GPIO_CLOCK_BASE, uint32BlockSize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return errors.New("mmap (CLOCK) failed")
	}
	i = (*[uint32BlockSize / 4]uint32)(unsafe.Pointer(&clk[0]))
	clkArr = i[:]
	//	The drive pads
	pads, err := syscall.Mmap(int(file.Fd()), GPIO_PADS, uint32BlockSize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return errors.New("mmap (PADS) failed")
	}
	i = (*[uint32BlockSize / 4]uint32)(unsafe.Pointer(&pads[0]))
	padsArry = i[:]
	return
}

/*


 digitalRead:
	Read the value of a given Pin, returning HIGH or LOW
 *********************************************************************************


int digitalRead (int pin)
{
  char c ;
  struct wiringPiNodeStruct *node = wiringPiNodes ;

  if ((pin & PI_GPIO_MASK) == 0)		// On-Board Pin
  {
    if (wiringPiMode == WPI_MODE_GPIO_SYS)	// Sys mode
    {
      if (sysFds [pin] == -1)
	return LOW ;

      lseek  (sysFds [pin], 0L, SEEK_SET) ;
      read   (sysFds [pin], &c, 1) ;
      return (c == '0') ? LOW : HIGH ;
    }
    else if (wiringPiMode == WPI_MODE_PINS)
      pin = pinToGpio [pin] ;
    else if (wiringPiMode == WPI_MODE_PHYS)
      pin = physToGpio [pin] ;
    else if (wiringPiMode != WPI_MODE_GPIO)
      return LOW ;

    if ((*(gpio + gpioToGPLEV [pin]) & (1 << (pin & 31))) != 0)
      return HIGH ;
    else
      return LOW ;
  }
  else
  {
    if ((node = wiringPiFindNode (pin)) == NULL)
      return LOW ;
    return node->digitalRead (node, pin) ;
  }
}
*/

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

func piBoardId() (pcbrev uint, bmodel uint, processor uint, manufacturer uint, ram uint, bWarranty uint, err error) {

	str := `Unable to determine boardinfo. If this is not a Raspberry Pi then you 
    are on your own as wiringPi is designed to support the 
    Raspberry Pi ONLY.\n`
	var ErrRevision error = errors.New(str)

	cpuinfo, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return 0, 0, 0, 0, 0, 0, ErrRevision

	}
	lines := strings.Split(string(cpuinfo), "\n")

	revisionValue := ""
	for _, line := range lines {
		fields := strings.Split(line, ":")
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		if key == "Revision" {
			ErrRevision = nil
			revisionValue = value
			break
		}
		//unicode.IsNumber
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
		return 0, 0, 0, 0, 0, 0, err
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

	return pcbrev, bmodel, processor, manufacturer, ram, bWarranty, ErrRevision

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
