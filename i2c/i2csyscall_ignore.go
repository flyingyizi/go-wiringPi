// +build ignore

package i2c

/*
#include <linux/i2c.h>
#include <linux/i2c-dev.h>
*/
import "C"

const (
	I2cSMBus      = C.I2C_SMBUS /* SMBus-level access */
	I2cSlaveForce = C.I2C_SLAVE_FORCE
	I2cSlave      = C.I2C_SLAVE

	I2cSMBusRead  = C.I2C_SMBUS_READ
	I2cSMBusWrite = C.I2C_SMBUS_WRITE

	I2cSMBusQuick        = C.I2C_SMBUS_QUICK
	I2cSMBusByte         = C.I2C_SMBUS_BYTE
	I2cSMBusByteData     = C.I2C_SMBUS_BYTE_DATA
	I2cSMBusWordData     = C.I2C_SMBUS_WORD_DATA
	I2cSMBusProcCall     = C.I2C_SMBUS_PROC_CALL
	I2cSMBusBlockData    = C.I2C_SMBUS_BLOCK_DATA
	I2cSMBusI2cBlockData = C.I2C_SMBUS_I2C_BLOCK_DATA
)

const (
	I2cSmBusBlockMax    = C.I2C_SMBUS_BLOCK_MAX     /* As specified in SMBus standard */
	I2cSmBusI2cBlockMax = C.I2C_SMBUS_I2C_BLOCK_MAX /* Not specified but we use same structure */
)

type i2c_smbus_ioctl_data C.struct_i2c_smbus_ioctl_data

const (
	Sizeofi2c_smbus_ioctl_data = C.sizeof_struct_i2c_smbus_ioctl_data
)
