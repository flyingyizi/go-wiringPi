// +build linux

package wiringPi

/*
#cgo linux, CFLAGS: -O2 -D_GNU_SOURCE -Wformat=2 -Wall -Wextra -Winline  -pipe -fPIC

#include "wiringPi/wiringPi/wiringSerial.c"
*/
import "C"
