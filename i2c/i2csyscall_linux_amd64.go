// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs i2csyscall_ignore.go

package i2c

const (
	I2cSMBus      = 0x720
	I2cSlaveForce = 0x706
	I2cSlave      = 0x703

	I2cSMBusRead  = 0x1
	I2cSMBusWrite = 0x0

	I2cSMBusQuick        = 0x0
	I2cSMBusByte         = 0x1
	I2cSMBusByteData     = 0x2
	I2cSMBusWordData     = 0x3
	I2cSMBusProcCall     = 0x4
	I2cSMBusBlockData    = 0x5
	I2cSMBusI2cBlockData = 0x8
)

const (
	I2cFuncs = 0x705
	I2cRDWR  = 0x707
)

const (
	I2cSmBusBlockMax    = 0x20
	I2cSmBusI2cBlockMax = I2cSmBusBlockMax
)

type i2c_smbus_ioctl_data struct {
	Write     uint8
	Command   uint8
	Pad_cgo_0 [2]byte
	Size      uint32
	Data      *[34]byte
}

const (
	Sizeofi2c_smbus_ioctl_data = 0x10
)

const (
	I2cFuncI2c = 0x1
)
