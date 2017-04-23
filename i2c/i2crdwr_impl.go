package i2c

import (
	"os"
	"syscall"
	"unsafe"
)

//1. IOCTL I2C_RDWR
//This method allows for simultaneous read/write and sending an uninterrupted
//sequence of message. Not all i2c devices support this method.
//Before performing i/o with this method, you should check whether the
//device supports this method using an ioctl I2C_FUNCS operation.
//Using this method, you do not need to perform an ioctl I2C_SLAVE
//operation -- it is done behind the scenes using the information embedded in the messages.

const (
	i2C_FUNCS = 0x0705 /* Get the adapter functionality */
	i2C_RDWR  = 0x0707 /* Combined R/W transfer (one stop only)*/
)

const (

	/* To determine what functionality is present */

	i2C_FUNC_I2C = 0x00000001

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

//https://github.com/ve3wwg/raspberry_pi/blob/master/mcp23017/i2c_funcs.c

func i2c_funcs_ioctl(f *os.File, data uintptr) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(f.Fd()), i2C_FUNCS, data); errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

func supportI2cRDWR(f *os.File) (b bool, err error) {
	var data uint64
	err = i2c_funcs_ioctl(f, uintptr(unsafe.Pointer(&data)))
	if err != nil {
		x := data & i2C_FUNC_I2C
		if x != 0 {
			b = true
		}
	}
	return
}

/*
struct i2c_rdwr_ioctl_data {

         structi2c_msg __user *msgs;        //指向i2c_msgs数组
         __u32nmsgs;                     // 消息的个数

};

msgs[] 数组成员包含了指向各自缓冲区的指针。这个函数会根据是否在消息中的flags置位I2C_M_RD来对缓
冲区进行读写。从机的地址以及是否使用10比特地址模式记录在每个消息中，忽略之前ioctl设置的结果。
*/

/*

int i2c_ioctl_write (int fd, uint8_t dev, uint8_t regaddr, uint16_t *data)
{
    int i, j = 0;
    int ret;
    uint8_t *buf;

    buf = malloc(1 + 2 * (sizeof(data) / sizeof(data[0])));
    if (buf == NULL) {
        return -ENOMEM;
    }

    buf[j ++] = regaddr;
    for (i = 0; i < (sizeof(data) / sizeof(data[0])); i ++) {
        buf[j ++] = (data[i] & 0xff00) >> 8;
        buf[j ++] = data[i] & 0xff;
    }

    struct i2c_msg messages[] = {
        {
            .addr = dev,
            .buf = buf,
            .len = sizeof(buf) / sizeof(buf[0]),
        },
    };

    struct i2c_rdwr_ioctl_data payload = {
        .msgs = messages,
        .nmsgs = sizeof(messages) / sizeof(messages[0]),
    };

    ret = ioctl(fd, I2C_RDWR, &payload);
    if (ret < 0) {
        ret = -errno;
    }

    free (buf);
    return ret;
}
*/
