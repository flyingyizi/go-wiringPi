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
	I2cFuncs = C.I2C_FUNCS
	I2cRDWR  = C.I2C_RDWR /* Combined R/W transfer (one stop only)*/
)

const (
	I2cSmBusBlockMax    = C.I2C_SMBUS_BLOCK_MAX /* As specified in SMBus standard */
	I2cSmBusI2cBlockMax = I2cSmBusBlockMax      /* Not specified but we use same structure */
)

type i2c_smbus_ioctl_data C.struct_i2c_smbus_ioctl_data

const (
	Sizeofi2c_smbus_ioctl_data = C.sizeof_struct_i2c_smbus_ioctl_data
)

const (

	/* To determine what functionality is present */

	I2cFuncI2c = C.I2C_FUNC_I2C // = 0x00000001

//#define I2C_FUNC_10BIT_ADDR		0x00000002
//#define I2C_FUNC_PROTOCOL_MANGLING	0x00000004 /* I2C_M_{REV_DIR_ADDR,NOSTART,..} */
//#define I2C_FUNC_SMBUS_PEC		0x00000008
//#define I2C_FUNC_SMBUS_BLOCK_PROC_CALL	0x00008000 /* SMBus 2.0 */
//#define I2C_FUNC_SMBUS_QUICK		0x00010000
//#define I2C_FUNC_SMBUS_READ_BYTE	0x00020000
//#define I2C_FUNC_SMBUS_WRITE_BYTE	0x00040000
//#define I2C_FUNC_SMBUS_READ_BYTE_DATA	0x00080000
//#define I2C_FUNC_SMBUS_WRITE_BYTE_DATA	0x00100000
//#define I2C_FUNC_SMBUS_READ_WORD_DATA	0x00200000
//#define I2C_FUNC_SMBUS_WRITE_WORD_DATA	0x00400000
//#define I2C_FUNC_SMBUS_PROC_CALL	0x00800000
//#define I2C_FUNC_SMBUS_READ_BLOCK_DATA	0x01000000
//#define I2C_FUNC_SMBUS_WRITE_BLOCK_DATA 0x02000000
//#define I2C_FUNC_SMBUS_READ_I2C_BLOCK	0x04000000 /* I2C-like block xfer  */
//#define I2C_FUNC_SMBUS_WRITE_I2C_BLOCK	0x08000000 /* w/ 1-byte reg. addr. */

)
