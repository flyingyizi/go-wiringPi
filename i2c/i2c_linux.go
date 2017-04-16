

package i2c

import (
	"os"
		"github.com/flyingyizi/go-wiringPi/board"
	"time"
	"unsafe"
)


/*
 * wiringPiI2CSetup:
 *	Open the I2C device, and regsiter the target device
 *********************************************************************************
 */

int I2CSetup (const int devId)
{
		info, _, err := board.GetBoardInfo()
		if err !=nil {
			return
		}
info.I2CDeviceName()


  return wiringPiI2CSetupInterface (device, devId) ;
}



/*
extern int wiringPiI2CRead           (int fd) ;
extern int wiringPiI2CReadReg8       (int fd, int reg) ;
extern int wiringPiI2CReadReg16      (int fd, int reg) ;

extern int wiringPiI2CWrite          (int fd, int data) ;
extern int wiringPiI2CWriteReg8      (int fd, int reg, int data) ;
extern int wiringPiI2CWriteReg16     (int fd, int reg, int data) ;

extern int wiringPiI2CSetupInterface (const char *device, int devId) ;
extern int wiringPiI2CSetup          (const int devId) ;


*/