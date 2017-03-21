// +build linux

package wiringPi

/*
#cgo linux, CFLAGS: -O2 -D_GNU_SOURCE -Wformat=2 -Wall -Wextra -Winline  -pipe -fPIC
#include "wiringPi/wiringPi/piHiPri.c"
#include "wiringPi/wiringPi/softPwm.c"
#include "wiringPi/wiringPi/wiringPi.c"

#include "wiringPi/wiringPi/wiringPiI2C.c"
#include "wiringPi/wiringPi/wiringPiSPI.c"
*/
import "C"
