package i2c

import (
	"os"
	"syscall"
	"unsafe"
)

//https://github.com/ve3wwg/raspberry_pi/blob/master/mcp23017/i2c_funcs.c

func i2c_funcs_ioctl(f *os.File, data uintptr) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(f.Fd()), I2cFuncs, data); errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

func supportI2cRDWR(f *os.File) (b bool, err error) {
	var data uint64
	err = i2c_funcs_ioctl(f, uintptr(unsafe.Pointer(&data)))
	if err != nil {
		x := data & I2cFuncI2c
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
