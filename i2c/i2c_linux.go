package i2c

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/flyingyizi/go-wiringPi/board"
)

const (
	i2cSLAVE  = 0x0703 // TODO(jbd): Allow users to use I2C_SLAVE_FORCE?
	i2cTENBIT = 0x0704
)

// Device represents an active connection to an I2C device.
type Device struct {
	f *os.File
}

// Read reads len(buf) bytes from the device.
func (d *Device) Read(buf []byte) error {
	return i2cTx(d.f, nil, buf)
}

// ReadReg is similar to Read but it reads from a register.
func (d *Device) ReadReg(reg byte, buf []byte) error {
	return i2cTx(d.f, []byte{reg}, buf)
}

// Write writes the buffer to the device. If it is required to write to a
// specific register, the register should be passed as the first byte in the
// given buffer.
func (d *Device) Write(buf []byte) (err error) {
	return i2cTx(d.f, buf, nil)
}

// WriteReg is similar to Write but writes to a register.
func (d *Device) WriteReg(reg byte, buf []byte) (err error) {
	// TODO(jbd): Do not allocate, not optimal.
	return i2cTx(d.f, append([]byte{reg}, buf...), nil)
}

// as a 10-bit address with TenBit.
// opens a 10-bit address example: d, err = i2c.Open( i2c.TenBit(0x78))
func Open(addr int) (d *Device, err error) {
	info, _, err := board.GetBoardInfo()
	if err != nil {
		return
	}
	device := info.I2CDeviceName()

	unmasked, tenbit := resolveAddr(addr)
	f, err := i2cOpen(device, unmasked, tenbit)

	return &(Device{f: f}), err
}

func (d *Device) Close() (err error) {
	err = i2cClose(d.f)
	return
}

const tenbitMask = 1 << 12

// TenBit marks an I2C address as a 10-bit address.
func TenBit(addr int) int {
	return addr | tenbitMask
}

// resolveAddr returns whether the addr is 10-bit masked or not.
// It also returns the unmasked address.
func resolveAddr(addr int) (unmasked int, tenbit bool) {
	return addr & (tenbitMask - 1), addr&tenbitMask == tenbitMask
}

// TODO(jbd): Support I2C_RETRIES and I2C_TIMEOUT at the driver and implementation level.
// Open opens a connection to an I2C device.
// All devices must be closed once they are no longer in use.
// For devices that use 10-bit I2C addresses, addr can be marked
func i2cOpen(device string, addr int, tenbit bool) (f *os.File, err error) {

	f, err = os.OpenFile(device, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return nil, err
	}

	if tenbit {
		if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), i2cTENBIT, uintptr(1)); errno != 0 {
			f.Close()
			//return syscall.Errno(errno)
			return nil, fmt.Errorf("cannot enable the 10-bit address mode on bus %v: %v", device, err)
		}
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), i2cSLAVE, uintptr(addr)); errno != 0 {
		f.Close()
		return nil, fmt.Errorf("error opening the address (%v) on the bus (%v): %v", addr, device, err)
	}
	return
}

func i2cClose(f *os.File) (err error) {
	if f != nil {
		err = f.Close()
	}
	return
}

// Tx first writes w (if not nil), then reads len(r)
// bytes from device into r (if not nil) in a single
// I2C transaction.
func i2cTx(f *os.File, w []byte, r []byte) error {
	if w != nil {
		if _, err := f.Write(w); err != nil {
			return err
		}
		f.Sync()
	}
	if r != nil {
		if _, err := io.ReadFull(f, r); err != nil {
			return err
		}
	}
	return nil
}

/*
 查看include/linux/i2c-dev.h文件，可以看到i2c-dev支持的IOCTL命令。
#define I2C_RETRIES                   0x0701                                   //设置收不到ACK时的重试次数
1．  设置重试次数
ioctl(fd, I2C_RETRIES,m); //这句话设置适配器收不到ACK时重试的次数为m。默认的重试次数为1。

#define I2C_TIMEOUT                0x0702                                   // 设置超时时限的jiffies
 ioctl(fd, I2C_TIMEOUT,m); //设置SMBus的超时时间为m，单位为jiffies。

#define I2C_SLAVE                      0x0703                                   设置从机地址
#define I2C_SLAVE_FORCE        0x0706                                  //强制设置从机地址
 ioctl(fd, I2C_SLAVE,addr);
ioctl(fd, #defineI2C_SLAVE_FORCE, addr);  //在调用read()和write()函数之前必须设置从机地址。这两行都可以设置从机的地址，区别是第二行无论内核中
                                          //是否已有驱动在使用这个地址都会成功，第一行则只在该地址空闲的情况下成功。由于i2c-dev创建的i2c_client不
										//加入i2c_adapter的client列表，所以不能防止其它线程使用同一地址，也不能防止驱动模块占用同一地址。


#define I2C_TENBIT                     0x0704                                   //选择地址位长:=0 for 7bit , != 0 for 10 bit
4．  设置地址模式
ioctl(file,I2C_TENBIT,select)     //如果select不等于0选择10比特地址模式，如果等于0选择7比特模式，默认7比特。只有适配器支持I2C_FUNC_10BIT_ADDR，这个请求才是有效的。

#define I2C_FUNCS                     0x0705                                   //获取适配器支持的功能
5．  获取适配器功能
ioctl(file,I2C_FUNCS,（unsignedlong *）funcs)  // 获取的适配器功能保存在funcs中。各比特的含义如
 // include/linux/i2c.h
#define I2C_FUNC_I2C                                                      0x00000001
#define I2C_FUNC_10BIT_ADDR                                    0x00000002
#define I2C_FUNC_PROTOCOL_MANGLING              0x00000004 //I2C_M_{REV_DIR_ADDR,NOSTART,..}
#define I2C_FUNC_SMBUS_PEC                                     0x00000008
#define I2C_FUNC_SMBUS_BLOCK_PROC_CALL     0x00008000  // SMBus 2.0
#define I2C_FUNC_SMBUS_QUICK                               0x00010000
#define I2C_FUNC_SMBUS_READ_BYTE                    0x00020000
#define I2C_FUNC_SMBUS_WRITE_BYTE                             0x00040000
#define I2C_FUNC_SMBUS_READ_BYTE_DATA        0x00080000
#define I2C_FUNC_SMBUS_WRITE_BYTE_DATA      0x00100000
#define I2C_FUNC_SMBUS_READ_WORD_DATA      0x00200000
#define I2C_FUNC_SMBUS_WRITE_WORD_DATA    0x00400000
#define I2C_FUNC_SMBUS_PROC_CALL                    0x00800000
#define I2C_FUNC_SMBUS_READ_BLOCK_DATA    0x01000000
#define I2C_FUNC_SMBUS_WRITE_BLOCK_DATA 0x02000000
#define I2C_FUNC_SMBUS_READ_I2C_BLOCK                  0x04000000  // I2C-like block xfer
#define I2C_FUNC_SMBUS_WRITE_I2C_BLOCK        0x08000000 // w/ 1-byte reg. addr.
#define I2C_FUNC_SMBUS_READ_I2C_BLOCK_2     0x10000000 // I2C-like block xfer
#define I2C_FUNC_SMBUS_WRITE_I2C_BLOCK_2   0x20000000 // w/ 2-byte reg. addr.

#define I2C_RDWR                       0x0707                                   // Combined R/W transfer (one STOP only)
6．  I2C层通信
ioctl(file,I2C_RDWR,(structi2c_rdwr_ioctl_data *)msgset); //这一行代码可以使用I2C协议和设备进行通信。它进行连续的读写，中间没有间歇。
                                                         //只有当适配器支持I2C_FUNC_I2C此命令才有效。参数是一个指针，指向一个结构体，它的定义如
														 struct i2c_rdwr_ioctl_data {
         structi2c_msg __user *msgs;  // 指向i2c_msgs数组
         __u32nmsgs;      //消息的个数

};
msgs[] 数组成员包含了指向各自缓冲区的指针。这个函数会根据是否在消息中的flags置位I2C_M_RD来对缓冲区进行读写。
从机的地址以及是否使用10比特地址模式记录在每个消息中，忽略之前ioctl设置的结果。

#define I2C_PEC                           0x0708                                   // != 0 to use PEC with SMBus
7．  设置SMBus PEC
ioctl(file,I2C_PEC,(long )select); //如果select不等于0选择SMBus PEC (packet error checking)，等于零则关闭这个功能，默认是关闭的。
这个命令只对SMBus传输有效。这个请求只在适配器支持I2C_FUNC_SMBUS_PEC时有效；如果不支持这个命令也是安全的，它不做任何工作。

#define I2C_SMBUS                     0x0720                                   /*SMBus transfer
8．  SMBus通信
ioctl(file, I2C_SMBUS, (i2c_smbus_ioctl_data*)msgset);  //这个函数和I2C_RDWR类似，参数的指针指向i2c_smbus_ioctl_data类型的变量，
它的定义如
struct i2c_smbus_ioctl_data {
         __u8read_write;
         __u8command;
         __u32size;
         unioni2c_smbus_data __user *data;
};

*/

/*
1.3     i2c_dev使用例程

要想在用户空间使用i2c适配器，首先要如3.1<!--[if gte mso 9]><![endif]-->节所示，选择某个适配器的设备节点打开，然后才能进行通信。
1.3.1   read()/write()

通信的方式有两种，一种是使用操作普通文件的接口read()和write()。这两个函数间接调用了i2c_master_recv和i2c_master_send。但是在使用之前需要使用I2C_SLAVE设置从机地址，设置可能失败，需要检查返回值。这种通信过程进行I2C层的通信，一次只能进行一个方向的传输。
下面的程序是ARM与E2PROM芯片通信的例子，如<!--[if supportFields]> REF _Ref283651035 /h <![endif]-->程序清单 3.5<!--[if gte mso 9]><![endif]--><!--[if supportFields]><![endif]-->所示。
程序清单 <!--[if supportFields]> STYLEREF 1 /s <![endif]-->3<!--[if supportFields]><![endif]-->.<!--[if supportFields]> SEQ 程序清单 /* ARABIC /s 1 <![endif]-->5<!--[if supportFields]><![endif]-->  使用read()/write()与i2c设备通信
#include <stdio.h>
#include <sys/ioctl.h>
#include <fcntl.h>
#include <linux/i2c-dev.h>
#include <linux/i2c.h>

#define CHIP                         "/dev/i2c-0"
#define CHIP_ADDR           0x50

int main()
{
         printf("hello,this is i2c tester/n");
         int fd =open(CHIP, O_RDWR);
         if (fd< 0) {
                   printf("open"CHIP"failed/n");
                   gotoexit;
         }

         if (ioctl(fd,I2C_SLAVE_FORCE, CHIP_ADDR) < 0) {            //设置芯片地址
                   printf("oictl:setslave address failed/n");
                   gotoclose;
         }

         struct                   i2c_msg msg;
         unsignedchar      rddata;
         unsignedchar      rdaddr[2] = {0, 0};                                         // 将要读取的数据在芯片中的偏移量
         unsignedchar      wrbuf[3] = {0, 0, 0x3c};                                  // 要写的数据，头两字节为偏移量

         printf("inputa char you want to write to E2PROM/n");
         wrbuf[2]= getchar();
         printf("writereturn:%d, write data:%x/n", write(fd, wrbuf, 3), wrbuf[2]);
         sleep(1);
         printf("writeaddress return: %d/n",write(fd, rdaddr, 2));       // 读取之前首先设置读取的偏移量
         printf("readdata return:%d/n", read(fd, &rddata, 1));
         printf("rddata:%c/n", rddata);
close:
         close(fd);
exit:
         return0;
}
1.3.2  I2C_RDWR

还可以使用I2C_RDWR实现同样的功能，如<!--[if supportFields]> REF _Ref283651333 /h <![endif]-->程序清单 3.6<!--[if gte mso 9]><![endif]--><!--[if supportFields]><![endif]-->所示。此时ioctl返回的值为执行成功的消息数。

程序清单 <!--[if supportFields]> STYLEREF 1 /s <![endif]-->3<!--[if supportFields]><![endif]-->.<!--[if supportFields]> SEQ 程序清单 /* ARABIC /s 1 <![endif]-->6<!--[if supportFields]><![endif]-->   使用I2C_RDWR与I2C设备通信

#include <stdio.h>
#include <sys/ioctl.h>
#include <fcntl.h>
#include <linux/i2c-dev.h>
#include <linux/i2c.h>

#define CHIP                         "/dev/i2c-0"
#define CHIP_ADDR           0x50

 int main()

{

         printf("hello,this is i2c tester/n");
         int fd =open(CHIP, O_RDWR);
         if (fd< 0) {
                   printf("open"CHIP"failed/n");

                   gotoexit;
         }

         struct                   i2c_msg msg;
         unsignedchar      rddata;
         unsignedchar      rdaddr[2] = {0, 0};
         unsignedchar      wrbuf[3] = {0, 0, 0x3c};
         printf("inputa char you want to write to E2PROM/n");
         wrbuf[2]= getchar();
         structi2c_rdwr_ioctl_data ioctl_data;

         structi2c_msg msgs[2];

         msgs[0].addr= CHIP_ADDR;
         msgs[0].len= 3;
         msgs[0].buf= wrbuf;
         ioctl_data.nmsgs= 1;
         ioctl_data.msgs= &msgs[0];

         printf("ioctlwrite,return :%d/n", ioctl(fd, I2C_RDWR, &ioctl_data));

         sleep(1);

         msgs[0].addr= CHIP_ADDR;
         msgs[0].len= 2;
         msgs[0].buf= rdaddr;
         msgs[1].addr= CHIP_ADDR;
         msgs[1].flags|= I2C_M_RD;
         msgs[1].len= 1;
         msgs[1].buf= &rddata;
         ioctl_data.nmsgs= 1;
         ioctl_data.msgs= msgs;

         printf("ioctlwrite address, return :%d/n", ioctl(fd, I2C_RDWR, &ioctl_data));
         ioctl_data.msgs= &msgs[1];
         printf("ioctlread, return :%d/n", ioctl(fd, I2C_RDWR, &ioctl_data));
         printf("rddata:%c/n", rddata);
close:
         close(fd);
exit:
         return0;
}

*/
