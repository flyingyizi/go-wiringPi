//+build linux

package wiringPi

/*
#include "wiringPi/wiringPi/wiringPi.h"
*/
import "C"

//Mode means wiringPi modes
type Mode int

const (
	PinsPinMode          Mode = C.WPI_MODE_PINS
	GPIOPinMode          Mode = C.WPI_MODE_GPIO     //GPIOPinMode Initialises the system into GPIO Pin mode and uses the	memory mapped hardware directly.
	GPIOSysPinMode       Mode = C.WPI_MODE_GPIO_SYS // Initialisation (again), however this time we are using the /sys/class/gpio  interface to the GPIO systems - slightly slower, but always usable as a non-root user, assuming the devices are already exported and setup correctly.
	PhysPinMode          Mode = C.WPI_MODE_PHYS     //PhysPinMode Initialises the system into Physical Pin mode and uses the 	memory mapped hardware directly
	PiFacePinMode        Mode = C.WPI_MODE_PIFACE
	UninitialisedPinMode Mode = C.WPI_MODE_UNINITIALISED
)

// PinMode corresponds to GPIO PIN mode.
type PinMode int

const (
	InputPinMode          PinMode = C.INPUT
	OutputPinMode         PinMode = C.OUTPUT
	PwmOutputPinMode      PinMode = C.PWM_OUTPUT
	GPIOClockPinMode      PinMode = C.GPIO_CLOCK
	SoftPwmOutputPinMode  PinMode = C.SOFT_PWM_OUTPUT
	SoftToneOutputPinMode PinMode = C.SOFT_TONE_OUTPUT
	PwmToneOutputPinMode  PinMode = C.PWM_TONE_OUTPUT
)

//LOW means logic low
const LOW int = 0

//HIGH means logic high
const HIGH int = 1

// PullDest : Pull up/down/none
type PullDest int

const (
	// PullOff : no pull up/down
	PullOff PullDest = C.PUD_OFF
	//PullDown : pull to ground
	PullDown PullDest = C.PUD_DOWN
	//PullUp : pull to 3.3V
	PullUp PullDest = C.PUD_UP
)

//PWM :Pulse Width Modulation
type PWM int

const (
	//PwmModeMS . The mark:space mode is traditional
	PwmModeMS PWM = C.PWM_MODE_MS
	// PwmModeBal . mode balanced
	PwmModeBal PWM = C.PWM_MODE_BAL
)

//InterruptLevel means Interrupt levels
type InterruptLevel int

const (
	IntEdgeSetup   InterruptLevel = C.INT_EDGE_SETUP
	IntEdgeFalling InterruptLevel = C.INT_EDGE_FALLING
	IntEdgeRising  InterruptLevel = C.INT_EDGE_RISING
	IntEdgeBoth    InterruptLevel = C.INT_EDGE_BOTH
)

/*
These functions work directly on the Raspberry Pi and also with
external GPIO modules such as GPIO expanders and so on, although
not all modules support all functions – e.g. the PiFace is
pre-configured for its fixed inputs and outputs, and the Raspberry Pi
has no on-board analog hardware.
*/

//pinMode sets the mode of a pin to either INPUT, OUTPUT,
//PWM_OUTPUT or GPIO_CLOCK. Note that only wiringPi
//pin 1 (BCM_GPIO 18) supports PWM output and only wiringPi
//pin 7 (BCM_GPIO 4) supports CLOCK output modes.
//
//This function has no effect when in Sys mode. If you
//need to change the pin mode, then you can do it with
//the gpio program in a script before you start your program.
func pinMode(pin int, mode PinMode) {

	C.pinMode(C.int(pin), C.int(mode))
}

//pullUpDnControl This sets the pull-up or pull-down resistor
//mode on the given pin, which should be set as an input.
//Unlike the Arduino, the BCM2835 has both pull-up an
//down internal resistors. The parameter pud should be;
//PUD_OFF, (no pull up/down), PUD_DOWN (pull to ground) or
//PUD_UP (pull to 3.3v) The internal pull up/down resistors
//have a value of approximately 50KΩ on the Raspberry Pi.
//This function has no effect on the Raspberry Pi’s GPIO pins
//when in Sys mode. If you need to activate a pull-up/pull-down,
//then you can do it with the gpio program in a script before
//you start your program.
func pullUpDnControl(pin int, pud PullDest) {

	C.pullUpDnControl(C.int(pin), C.int(pud))
}

//DigitalWrite Writes the value HIGH or LOW (1 or 0) to the
//given pin which must have been previously set as an output.
//WiringPi treats any non-zero number as HIGH,
//however 0 is the only representation of LOW.
func DigitalWrite(pin int, value int) {

	C.digitalWrite(C.int(pin), C.int(value))
}

//pwmWrite Writes the value to the PWM register for the given pin.
//The Raspberry Pi has one on-board PWM pin, pin 1
//(BMC_GPIO 18, Phys 12) and the range is 0-1024. Other PWM
//devices may have other PWM ranges.
//This function is not able to control the Pi’s on-board
//PWM when in Sys mode.
func pwmWrite(pin int, value int) {

	C.pwmWrite(C.int(pin), C.int(value))
}

//digitalRead function returns the value read at the given pin.
//It will be HIGH or LOW (1 or 0) depending on the logic level at the pin.
func digitalRead(pin int) int {

	ret := int(C.digitalRead(C.int(pin)))
	return ret
}

//analogRead returns the value read on the supplied analog input pin.
//You will need to register additional analog modules
//to enable this function for devices such as the Gertboard,
//quick2Wire analog board, etc.
func analogRead(pin int) int {

	ret := int(C.analogRead(C.int(pin)))
	return ret
}

//AnalogWrite writes the given value to the supplied
//analog pin. You will need to register additional analog
//modules to enable this function for devices such as the Gertboard.
func AnalogWrite(pin int, value int) {
	C.analogWrite(C.int(pin), C.int(value))
}

//Setup Must be called once at the start of your program execution.
//* Default setup: Initialises the system into wiringPi Pin mode and uses the
//*	memory mapped hardware directly.
//* Changed now to revert to "gpio" mode if we're running on a Compute Module.
func Setup() int {
	ret := C.wiringPiSetup()
	return int(ret)
}

// SetPinMode : Sets the mode of a pin to be input, output or PWM output
func SetPinMode(pin int, mode PinMode) {
	C.pinMode(C.int(pin), C.int(mode))
}
