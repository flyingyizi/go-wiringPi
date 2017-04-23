package i2c

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

const (
	/* this is for i2c-dev.c	*/
	i2cSLAVE = 0x0703 /* Change slave address			*/
	/* Attn.: Slave address is 7 or 10 bits */
	i2cSLAVEForce = 0x0706 /* Change slave address			*/
	/* Attn.: Slave address is 7 or 10 bits */
	/* This changes the address, even if it */
	/* is already taken!			*/
	i2cTENBIT = 0x0704 /* I2C_TENBIT:0 for 7 bit addrs, != 0 for 10 bit	*/
)

// Device represents an active connection to an I2C device.
type Device struct {
	sync.Mutex

	// File used to represent the bus once it's opened
	f *os.File
	// example:"/dev/i2c-2"
	name string

	masterIsBigEndian bool // if BigEndian it is true, else false
}

// Open opens a connection to an I2C slave device.
// All devices must be closed once they are no longer in use.
// TODO(jbd): Support I2C_RETRIES and I2C_TIMEOUT at the driver and implementation level.
func Open(device string) (d *Device, err error) {
	f, err := os.OpenFile(device, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return nil, err
	}

	m := getEndian()

	return &(Device{f: f, name: device, masterIsBigEndian: m}), err
}

func (d *Device) Close() (err error) {
	if d != nil {
		err = d.f.Close()
	}
	return
}

const tenbitMask = 1 << 12

// TenBit marks an I2C address as a 10-bit address.
func TenBit(addr int) int {
	return addr | tenbitMask
}

// SetAddr set the I2C slave address for all subsequent I2C device transfers
// For devices that use 10-bit I2C addresses, addr can be marked
// as a 10-bit address with TenBit.
// set a 10-bit address example: err = d.SetAddr( i2c.TenBit(0x78))
func (d *Device) SetAddr(addr int) (err error) {
	unmasked := addr & (tenbitMask - 1)     //get the unmasked address
	tenbit := addr&tenbitMask == tenbitMask //whether the addr is 10-bit masked or not

	d.Lock()
	defer d.Unlock()

	if tenbit {
		if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, d.f.Fd(), i2cTENBIT, uintptr(1)); errno != 0 {
			d.f.Close()
			return fmt.Errorf("cannot enable the 10-bit address mode : %v", errno)
		}
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, d.f.Fd(), i2cSLAVE, uintptr(unmasked)); errno != 0 {
		d.f.Close()
		return fmt.Errorf("error opening the address (%v) on bus (%s) : %v", addr, d.name, errno)
	}
	return
}

func (d *Device) GetName() string {
	return d.name
}

//SmbusWriteQuick 	Sends a single bit to the device (in place of the Rd/Wr bit shown in Listing 8.1).
func (d *Device) SmbusWriteQuick(value uint8) error {
	d.Lock()
	defer d.Unlock()

	return i2c_smbus_write_quick(d.f, value)
}

//SmbusReadByte   Reads a single byte from the device without specifying a location offset.
//Uses the same offset as the previously issued command.
func (d *Device) SmbusReadByte() (data uint8, err error) {
	d.Lock()
	defer d.Unlock()

	data, err = i2c_smbus_read_byte(d.f)
	return
}

//SmbusWriteByte  	Sends a single byte to the device at the same memory offset as the previously issued command.
func (d *Device) SmbusWriteByte(value uint8) error {
	d.Lock()
	defer d.Unlock()

	return i2c_smbus_write_byte(d.f, value)
}

//SmbusReadByteData   	Reads a single byte from the device at a specified offset.
func (d *Device) SmbusReadByteData(command uint8) (data uint8, err error) {
	d.Lock()
	defer d.Unlock()

	data, err = i2c_smbus_read_byte_data(d.f, command)
	return
}

//SmbusWriteByteData    Sends a single byte to the device at a specified offset.
func (d *Device) SmbusWriteByteData(command uint8, value uint8) (err error) {
	d.Lock()
	defer d.Unlock()

	return i2c_smbus_write_byte_data(d.f, command, value)
}

//SmbusReadWordData   	Reads 2 bytes from the specified offset.
func (d *Device) SmbusReadWordData(command uint8) (data uint16, err error) {
	d.Lock()
	defer d.Unlock()

	return i2c_smbus_read_word_data(d.f, command)
}

//SmbusWriteWordData    	Sends 2 bytes to the specified offset.
func (d *Device) SmbusWriteWordData(command uint8, value uint16) (err error) {
	d.Lock()
	defer d.Unlock()

	return i2c_smbus_write_word_data(d.f, command, value)
}

func (d *Device) SmbusProcessCall(command uint8, value uint16) (data uint16, err error) {
	d.Lock()
	defer d.Unlock()

	data, err = i2c_smbus_process_call(d.f, command, value)
	return
}

//SmbusReadBlockData   	Reads a block of data from the specified offset.
func (d *Device) SmbusReadBlockData(command uint8) (block []byte, err error) {
	d.Lock()
	defer d.Unlock()

	block, err = i2c_smbus_read_block_data(d.f, command)
	return
}

//SmbusWriteBlockData   	Sends a block of data (<= 32 bytes) to the specified offset.
func (d *Device) SmbusWriteBlockData(command uint8, length uint8, value []byte) (err error) {
	d.Lock()
	defer d.Unlock()

	return i2c_smbus_write_block_data(d.f, command, length, value)
}

func (d *Device) SmbusWriteI2cBlockData(command uint8, length uint8, value []byte) (err error) {
	d.Lock()
	defer d.Unlock()

	return i2c_smbus_write_i2c_block_data(d.f, command, length, value)
}

// SysfsRead reads len(buf) bytes from the device.
func (d *Device) SysfsRead(buf []byte) error {
	return i2cTx(d.f, nil, buf)
}

// SysfsReadReg is similar to Read but it reads from a register.
func (d *Device) SysfsReadReg(reg byte, buf []byte) error {
	return i2cTx(d.f, []byte{reg}, buf)
}

// SysfsWrite writes the buffer to the device. If it is required to write to a
// specific register, the register should be passed as the first byte in the
// given buffer.
func (d *Device) SysfsWrite(buf []byte) (err error) {
	return i2cTx(d.f, buf, nil)
}

// SysfsWriteReg is similar to Write but writes to a register.
func (d *Device) SysfsWriteReg(reg byte, buf []byte) (err error) {
	// TODO(jbd): Do not allocate, not optimal.
	return i2cTx(d.f, append([]byte{reg}, buf...), nil)
}

// ref https://github.com/virtao/GoEndian/blob/master/endian.go
//true = big endian, false = little endian
func getEndian() (ret bool) {
	//以下代码判断机器大小端
	const intsize int = int(unsafe.Sizeof(0))
	var i int = 0x1
	bs := (*[intsize]byte)(unsafe.Pointer(&i))
	if bs[0] == 0 {
		return true
	}
	return false
}
