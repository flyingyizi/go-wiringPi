// +build linux

package wiringPi

/*
#cgo linux, CFLAGS: -O2 -D_GNU_SOURCE -Wformat=2 -Wall -Wextra -Winline -I. -pipe -fPIC
#include "wiringPi/wiringPi/softTone.c"
*/
import "C"
