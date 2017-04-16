// +build linux

package serial

import (
	"os"
	"syscall"
	"time"
	"unsafe"
)

