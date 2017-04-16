package spi

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

/*
#cgo linux, CFLAGS: -O2 -D_GNU_SOURCE -Wformat=2 -Wall -Wextra -Winline  -pipe -fPIC
#include <sys/ioctl.h>
#include <asm/ioctl.h>
#include <linux/spi/spidev.h>
*/
import "C"

//Mode means wiringPi modes
type Mode int

const (
	SpiIOCWRMode        = C.SPI_IOC_WR_MODE
	spiIOCWRMaxSpeedHz  = C.SPI_IOC_WR_MAX_SPEED_HZ
	spiIOCWRBitsPerWord = C.SPI_IOC_WR_BITS_PER_WORD
	defaultDelayms      = 0
	defaultSPIBPW       = 8
	defaultSPISpeed     = 1000000

	spiIOCMessage0    = 1073769216 //0x40006B00
	spiIOCIncrementor = 2097152    //0x200000

)

/*
const (
	spiIOCWrMode        = 0x40016B01
	spiIOCWrBitsPerWord = 0x40016B03
	spiIOCWrMaxSpeedHz  = 0x40046B04

	spiIOCRdBitsPerWord = 0x80016B03
	spiIOCRdMaxSpeedHz  = 0x80046B04
	spiIOCMessage0    = 1073769216 //0x40006B00
	spiIOCIncrementor = 2097152    //0x200000


)
*/

// The SPI bus parameters
//	Variables as they need to be passed as pointers later on

const (
	spiDev0 string = "/dev/spidev0.0"
	spiDev1 string = "/dev/spidev0.1"
)

/*
 * wiringPiSPISetup:
 *	Open the SPI device, and set it up, etc. in the default MODE 0
 *********************************************************************************
 */
func spiOpen(channel byte, model uint8, speed uint32, bpw uint8) (f *os.File, err error) {

	model &= 3   // Mode is 0, 1, 2 or 3
	channel &= 1 // Channel is 0 or 1
	var device string
	if channel == 0 {
		device = spiDev0
	} else {
		device = spiDev1
	}

	f, err = os.OpenFile(device, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return nil, err //"Unable to open SPI device: %s\n"
	}
	//glog.V(3).Infof("spi: sucessfully opened file /dev/spidev0.%v", channel)

	if err = spiSetMode(f, model); err != nil {
		f.Close()
		return nil, err
	}
	if _, err = spiSetSpeed(f, speed); err != nil {
		f.Close()
		return nil, err
	}
	if _, err = spiSetBPW(f, bpw); err != nil {
		f.Close()
		return nil, err
	}

	return
}

func spiSetMode(f *os.File, mode uint8) error {
	fmt.Printf("spi: setting spi mode to %v", mode)

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), SpiIOCWRMode, uintptr(unsafe.Pointer(&mode)))
	if errno != 0 {
		err := syscall.Errno(errno)
		return err
	}
	fmt.Printf("spi: mode set to %v", mode)
	return nil
}

func spiSetSpeed(f *os.File, speed uint32) (speedHz uint32, err error) {
	if speed <= 0 {
		speed = defaultSPISpeed
	}

	//glog.V(3).Infof("spi: setting spi speedMax to %v", speed)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), spiIOCWRMaxSpeedHz, uintptr(unsafe.Pointer(&speed)))
	if errno != 0 {
		err = syscall.Errno(errno)
		//glog.V(3).Infof("spi: failed to set speedMax due to %v", err.Error())
		return
	}
	//glog.V(3).Infof("spi: speedMax set to %v", speed)
	speedHz = speed
	return
}

func spiSetBPW(f *os.File, bpw uint8) (bitsPerWord uint8, err error) {
	if bpw <= 0 {
		bpw = defaultSPIBPW
	}

	//glog.V(3).Infof("spi: setting spi bpw to %v", bpw)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), spiIOCWRBitsPerWord, uintptr(unsafe.Pointer(&bpw)))
	if errno != 0 {
		err = syscall.Errno(errno)
		//glog.V(3).Infof("spi: failed to set bpw due to %v", err.Error())
		return
	}
	//glog.V(3).Infof("spi: bpw set to %v", bpw)
	bitsPerWord = bpw
	return
}

func spiIOCMessageN(n uint32) uint32 {
	return (spiIOCMessage0 + (n * spiIOCIncrementor))
	//i := uint32(C.SPI_IOC_MESSAGE(C.__u32(n)))
	//return uint32(i)
}

/*
 * spiTx:
 *	Write and Read a block of data over the SPI bus.
 *	Note the data ia being read into the transmit buffer, so will
 *	overwrite it!
 *	This is also a full-duplex operation.
 *********************************************************************************
 */
//https://github.com/kidoman/embd/blob/master/host/generic/spibus.go
func spiTx(f *os.File, dataBuffer []uint8, speed uint32, bpw uint8) error {

	if f == nil {
		return errors.New("bad")
	}

	len := len(dataBuffer)

	// struct  spi_ioc_transfer
	var dataCarrier C.struct_spi_ioc_transfer //= C.struct_spi_ioc_transfer{id, 21}

	dataCarrier.len = C.__u32(len)
	dataCarrier.tx_buf = C.__u64(uintptr(unsafe.Pointer(&dataBuffer[0])))
	dataCarrier.rx_buf = C.__u64(uintptr(unsafe.Pointer(&dataBuffer[0])))
	dataCarrier.speed_hz = C.__u32(speed)
	dataCarrier.bits_per_word = C.__u8(bpw)
	//spi.tx_buf        = (unsigned long)data ;
	//spi.rx_buf        = (unsigned long)data ;
	//spi.len           = len ;
	//spi.delay_usecs   = spiDelay ;
	//spi.speed_hz      = spiSpeeds [channel] ;
	//spi.bits_per_word = spiBPW ;

	//glog.V(3).Infof("spi: sending dataBuffer %v with carrier %v", dataBuffer, dataCarrier)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(spiIOCMessageN(1)), uintptr(unsafe.Pointer(&dataCarrier)))
	if errno != 0 {
		err := syscall.Errno(errno)
		//glog.V(3).Infof("spi: failed to read due to %v", err.Error())
		return err
	}
	//glog.V(3).Infof("spi: read into dataBuffer %v", dataBuffer)
	return nil

	// Mentioned in spidev.h but not used in the original kernel documentation
	//	test program )-:

	//memset (&spi, 0, sizeof (spi)) ;
	//
	//spi.tx_buf        = (unsigned long)data ;
	//spi.rx_buf        = (unsigned long)data ;
	//spi.len           = len ;
	//spi.delay_usecs   = spiDelay ;
	//spi.speed_hz      = spiSpeeds [channel] ;
	//spi.bits_per_word = spiBPW ;
	/*
	   每个 spi_ioc_transfer都可以包含读和写的请求，其中读和写的长度必须相等。所以成员len不是
	   tx_buf和rx_buf缓冲的长度之和，而是它们各自的长度。SPI控制器驱动会先将tx_buf写到SPI总线上，
	   然后再读取len长度的内容到rx_buf。如果只想进行一个方向的传输，把另一个方向的缓冲置为0就可以了。
	   speed_hz和bits_per_word这两个成员可以为每次通信配置不同的通信速率（必须小于spi_device
	   的max_speed_hz）和字长，如果它们为0的话就会使用spi_device中的配置。
	   delay_usecs可以指定两个spi_ioc_transfer之间的延时，单位是微妙。一般不用定义。
	   cs_change指定这个cs_change结束之后是否需要改变片选线。一般针对同一设备的连续的几个
	   spi_ioc_transfer，只有最后一个需要将这个成员置位。这样省去了来回改变片选线的时间，有助于提高通信速率。
	   struct spi_ioc_transfer {
	   __u64 tx_buf; // 写数据缓冲
	   __u64 rx_buf; // 读数据缓冲
	   __u32 len; // 缓冲的长度
	   __u32 speed_hz; // 通信的时钟频率
	   __u16 delay_usecs; // 两个spi_ioc_transfer之间的延时
	   __u8 bits_per_word; // 字长（比特数）
	   __u8 cs_change; // 是否改变片选
	   __u32 pad;
	   };
	*/
	/*
	   type spiIOCTransfer struct {
	   	txBuf uint64
	   	rxBuf uint64
	   	length      uint32
	   	speedHz     uint32
	   	delayus     uint16
	   	bitsPerWord uint8
	   	csChange    uint8
	   	pad         uint32
	   }
	*/

}

// Device represents an active connection to an I2C device.
type Device struct {
	f *os.File

	// struct  spi_ioc_transfer
	spiTransferData C.struct_spi_ioc_transfer //= C.struct_spi_ioc_transfer{id, 21}

	channel byte
	mode    byte
	speed   uint32
	bpw     uint8
	delayms int

	mu sync.Mutex
}


func (d *Device) Open(channel byte, model uint8, speed uint32, bpw uint8)  error {
f, err :=spiOpen(channel byte, model uint8, speed uint32, bpw uint8) (f *os.File, err error) 
if err !=nil {
	return
}
d.f =f
d.channel=channel
d.mode=model
d.speed=speed
return

}


// TransferAndReceiveData, fistly sending databuffer, then read into data buffer
func (d *Device) Tx(dataBuffer []uint8) (err error) {

    err =spiTx(d.f, dataBuffer, d.speed, d.bpw) 
return
}


func (d *Device) Write(data []byte) (n int, err error) {
	return d.f.Write(data)
}


func (d *Device) ReceiveData(len int) ([]uint8, error) {
	data := make([]uint8, len)
	err := d.Tx(data)
	if  err != nil {
		return nil, err
	}
	return data, err
}

func (d *Device) TransferAndReceiveByte(data byte) (byte, error) {

	dt := [1]uint8{uint8(data)}
	if err := d.Tx(dt[:]); err != nil {
		return 0, err
	}
	return d[0], nil
}

func (d *Device) ReceiveByte() (byte, error) {

	var dt [1]uint8
	if err := d.Tx(dt[:]); err != nil {
		return 0, err
	}
	return byte(d[0]), nil
}


func (d *Device) Close() error {
	b.mu.Lock()
	defer d.mu.Unlock()

	return d.f.Close()
}