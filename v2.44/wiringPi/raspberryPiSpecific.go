package wiringPi

/*
#include "wiringPi/wiringPi/wiringPi.h"
*/
import "C"

//These functions are not part of the core wiringPi set,
//but act specifically on the Raspberry Pi hardware itself.
//Some external hardware driver modules may provide some of
//this functionality though.

//digitalWriteByte writes the 8-bit byte supplied to the first
//8 GPIO pins. It’s the fastest way to set all 8 bits
//at once to a particular value, although it still
//takes two write operations to the Pi’s GPIO hardware.
func digitalWriteByte(value int) {

	C.digitalWriteByte(C.int(value))
}

/*

    pwmSetMode (int mode) ;
The PWM generator can run in 2 modes – “balanced”
and “mark:space”. The mark:space mode is traditional,
however the default mode in the Pi is “balanced”.
You can switch modes by supplying the parameter: PWM_MODE_BAL or PWM_MODE_MS.

    pwmSetRange (unsigned int range) ;
This sets the range register in the PWM generator.
The default is 1024.

    pwmSetClock (int divisor) ;
This sets the divisor for the PWM clock.
Note: The PWM control functions can not be used when
in Sys mode. To understand more about the PWM system,
you’ll need to read the Broadcom ARM peripherals manual.

    piBoardRev (void) ;
This returns the board revision of the Raspberry Pi.
It will be either 1 or 2. Some of the BCM_GPIO pins
changed number and function when moving from board
revision 1 to 2, so if you are using BCM_GPIO pin numbers,
then you need to be aware of the differences.

    wpiPinToGpio (int wPiPin) ;

This returns the BCM_GPIO pin number of the supplied
wiringPi pin. It takes the board revision into account.

    physPinToGpio (int physPin) ;
This returns the BCM_GPIO pin number of the supplied physical
pin on the P1 connector.

    setPadDrive (int group, int value) ;
This sets the “strength” of the pad drivers for
a particular group of pins. There are 3 groups of
pins and the drive strength is from 0 to 7. Do not use
this unless you know what you are doing.
*/
