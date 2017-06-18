package i2c

//2. IOCTL SMBUS
//This method of i/o is more powerful but the resulting code is more verbose.
//This method can be used if the device does not support the I2C_RDWR method.
//Using this method, you do need to perform an ioctl I2C_SLAVE operation (or, if
//the device is busy, an I2C_SLAVE_FORCE operation).

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func i2c_smbus_ioctl(f *os.File, data uintptr) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(f.Fd()), I2cSMBus, data); errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

func i2c_smbus_access(f *os.File, read_write uint8, command uint8, size uint32, data interface{}) (err error) {

	/*static inline __s32 i2c_smbus_access(int file, char read_write, __u8 command,
	                                       int size, union i2c_smbus_data *data)
	  {
	  	struct i2c_smbus_ioctl_data args;

	  	args.read_write = read_write;
	  	args.command = command;
	  	args.size = size;
	  	args.data = data;
	  	return ioctl(file,I2C_SMBUS,&args);
	  }
	*/
	args := i2c_smbus_ioctl_data{Write: read_write, Command: command, Size: size}
	if data == nil {
		args.Data = (*[34]byte)(unsafe.Pointer(nil))
		err = i2c_smbus_ioctl(f, uintptr(unsafe.Pointer(&args)))

	} else {
		switch data.(type) {
		case *uint8:
			x, ok := data.(*uint8)
			if ok {
				args.Data = (*[34]byte)(unsafe.Pointer(x))
				err = i2c_smbus_ioctl(f, uintptr(unsafe.Pointer(&args)))
			}
		case *uint16:
			x, ok := data.(*uint16)
			if ok {
				args.Data = (*[34]byte)(unsafe.Pointer(x))
				err = i2c_smbus_ioctl(f, uintptr(unsafe.Pointer(&args)))
			}
		case []byte:
			x, ok := data.([]byte)
			if ok {
				args.Data = (*[34]byte)(unsafe.Pointer(&x[0]))
				err = i2c_smbus_ioctl(f, uintptr(unsafe.Pointer(&args)))
			}
		default:
			err = fmt.Errorf("error do i2c_smbus_ioctl_data")
		}
	}
	return
}

//i2c_smbus_write_quick()	Sends a single bit to the device (in place of the Rd/Wr bit shown in Listing 8.1).
func i2c_smbus_write_quick(f *os.File, value uint8) (err error) {
	/*static inline __s32 i2c_smbus_write_quick(int file, __u8 value)
	  {
	  	return i2c_smbus_access(file,value,0,I2C_SMBUS_QUICK,NULL);
	  }
	*/
	err = i2c_smbus_access(f, value /*read_write*/, 0 /*command*/, I2cSMBusQuick /*size*/, nil /*data*/)
	return
}

//i2c_smbus_read_byte()	Reads a single byte from the device without specifying a
//location offset. Uses the same offset as the previously issued command.
func i2c_smbus_read_byte(f *os.File) (data uint8, err error) {
	/*static inline __s32 i2c_smbus_read_byte(int file)
	  {
	  	union i2c_smbus_data data;
	  	if (i2c_smbus_access(file,I2C_SMBUS_READ,0,I2C_SMBUS_BYTE,&data))
	  		return -1;
	  	else
	  		return 0x0FF & data.byte;
	  }*/
	err = i2c_smbus_access(f, I2cSMBusRead /*read_write*/, 0 /*command*/, I2cSMBusByte /*size*/, &data /*data*/)
	return
}

//i2c_smbus_write_byte()	Sends a single byte to the device at the same memory
//offset as the previously issued command.
func i2c_smbus_write_byte(f *os.File, value uint8) (err error) {
	/*static inline __s32 i2c_smbus_write_byte(int file, __u8 value)
	  {
	  	return i2c_smbus_access(file,I2C_SMBUS_WRITE,value,
	  	                        I2C_SMBUS_BYTE,NULL);
	  }*/
	err = i2c_smbus_access(f, I2cSMBusWrite /*read_write*/, value /*command*/, I2cSMBusByte /*size*/, nil /*data*/)
	return
}

//i2c_smbus_read_byte_data()	Reads a single byte from the device at a specified offset.
func i2c_smbus_read_byte_data(f *os.File, command uint8) (data uint8, err error) {
	/*static inline __s32 i2c_smbus_read_byte_data(int file, __u8 command)
	  {
	  	union i2c_smbus_data data;
	  	if (i2c_smbus_access(file,I2C_SMBUS_READ,command,
	  	                     I2C_SMBUS_BYTE_DATA,&data))
	  		return -1;
	  	else
	  		return 0x0FF & data.byte;
	  }*/
	err = i2c_smbus_access(f, I2cSMBusRead /*read_write*/, command /*command*/, I2cSMBusByteData /*size*/, &data /*data*/)
	return
}

//i2c_smbus_write_byte_data()	Sends a single byte to the device at a specified offset.
func i2c_smbus_write_byte_data(f *os.File, command uint8, value uint8) (err error) {
	/*static inline __s32 i2c_smbus_write_byte_data(int file, __u8 command,
	                                                __u8 value)
	  {
	  	union i2c_smbus_data data;
	  	data.byte = value;
	  	return i2c_smbus_access(file,I2C_SMBUS_WRITE,command,
	  	                        I2C_SMBUS_BYTE_DATA, &data);
	  }*/
	err = i2c_smbus_access(f, I2cSMBusWrite /*read_write*/, command /*command*/, I2cSMBusByteData /*size*/, &value /*data*/)
	return
}

//i2c_smbus_read_word_data()	Reads 2 bytes from the specified offset.
func i2c_smbus_read_word_data(f *os.File, command uint8) (data uint16, err error) {
	/*static inline __s32 i2c_smbus_read_word_data(int file, __u8 command)
	  {
	  	union i2c_smbus_data data;
	  	if (i2c_smbus_access(file,I2C_SMBUS_READ,command,
	  	                     I2C_SMBUS_WORD_DATA,&data))
	  		return -1;
	  	else
	  		return 0x0FFFF & data.word;
	  }*/
	err = i2c_smbus_access(f, I2cSMBusRead /*read_write*/, command /*command*/, I2cSMBusWordData /*size*/, &data /*data*/)
	return
}

//i2c_smbus_write_word_data()	Sends 2 bytes to the specified offset.
func i2c_smbus_write_word_data(f *os.File, command uint8, value uint16) (err error) {
	/*static inline __s32 i2c_smbus_write_word_data(int file, __u8 command,
	                                                __u16 value)
	  {
	  	union i2c_smbus_data data;
	  	data.word = value;
	  	return i2c_smbus_access(file,I2C_SMBUS_WRITE,command,
	  	                        I2C_SMBUS_WORD_DATA, &data);
	  }*/
	err = i2c_smbus_access(f, I2cSMBusWrite /*read_write*/, command /*command*/, I2cSMBusWordData /*size*/, &value /*data*/)
	return
}

func i2c_smbus_process_call(f *os.File, command uint8, value uint16) (data uint16, err error) {
	/*static inline __s32 i2c_smbus_process_call(int file, __u8 command, __u16 value)
	  {
	  	union i2c_smbus_data data;
	  	data.word = value;
	  	if (i2c_smbus_access(file,I2C_SMBUS_WRITE,command,
	  	                     I2C_SMBUS_PROC_CALL,&data))
	  		return -1;
	  	else
	  		return 0x0FFFF & data.word;
	  }*/
	data = value
	err = i2c_smbus_access(f, I2cSMBusWrite /*read_write*/, command /*command*/, I2cSMBusProcCall /*size*/, &data /*data*/)
	return
}

/* Returns the read bytes */
//i2c_smbus_read_block_data()	Reads a block of data from the specified offset.
func i2c_smbus_read_block_data(f *os.File, command uint8) ([]byte, error) {
	/*static inline __s32 i2c_smbus_read_block_data(int file, __u8 command,
	                                                __u8 *values)
	  {
	  	union i2c_smbus_data data;
	  	int i;
	  	if (i2c_smbus_access(file,I2C_SMBUS_READ,command,
	  	                     I2C_SMBUS_BLOCK_DATA,&data))
	  		return -1;
	  	else {
	  		for (i = 1; i <= data.block[0]; i++)
	  			values[i-1] = data.block[i];
	  			return data.block[0];
	  	}
	  }*/
	block := make([]byte, I2cSmBusBlockMax+2, I2cSmBusBlockMax+2)
	err := i2c_smbus_access(f, I2cSMBusRead /*read_write*/, command /*command*/, I2cSMBusBlockData /*size*/, block /*data*/)
	len := len(block)
	if (len > 0) && err == nil {
		return block[1 : 1+len], nil
	}
	return block, fmt.Errorf("i2c_smbus_read_block_data: can not read ")
}

//i2c_smbus_write_block_data()	Sends a block of data (<= 32 bytes) to the specified offset.
func i2c_smbus_write_block_data(f *os.File, command uint8, length uint8, value []byte) (err error) {
	/*static inline __s32 i2c_smbus_write_block_data(int file, __u8 command,
	                                                 __u8 length, __u8 *values)
	  {
	  	union i2c_smbus_data data;
	  	int i;
	  	if (length > 32)
	  		length = 32;
	  	for (i = 1; i <= length; i++)
	  		data.block[i] = values[i-1];
	  	data.block[0] = length;
	  	return i2c_smbus_access(file,I2C_SMBUS_WRITE,command,
	  	                        I2C_SMBUS_BLOCK_DATA, &data);
	  }*/
	if length > 32 {
		length = 32
	}
	value = value[:length]
	value = append([]byte{length}, value[0:]...)

	err = i2c_smbus_access(f, I2cSMBusWrite /*read_write*/, command /*command*/, I2cSMBusBlockData /*size*/, value /*data*/)
	return
}

func i2c_smbus_write_i2c_block_data(f *os.File, command uint8, length uint8, value []byte) (err error) {
	/*static inline __s32 i2c_smbus_write_i2c_block_data(int file, __u8 command,
	                                                 __u8 length, __u8 *values)
	  {
	  	union i2c_smbus_data data;
	  	int i;
	  	if (length > 32)
	  		length = 32;
	  	for (i = 1; i <= length; i++)
	  		data.block[i] = values[i-1];
	  	data.block[0] = length;
	  	return i2c_smbus_access(file,I2C_SMBUS_WRITE,command,
	  	                        I2C_SMBUS_I2C_BLOCK_DATA, &data);
	  }*/
	if length > 32 {
		length = 32
	}
	value = value[:length]
	value = append([]byte{length}, value[0:]...)

	err = i2c_smbus_access(f, I2cSMBusWrite /*read_write*/, command /*command*/, I2cSMBusI2cBlockData /*size*/, value /*data*/)
	return
}

/*
Suppose we are on a 32-bit machine.
If it is little endian, the x in the memory will be something like:

       higher memory
          ----->
    +----+----+----+----+
    |0x01|0x00|0x00|0x00|
    +----+----+----+----+
    A
    |
   &x

so (char*)(*x) == 1, and *y+48 == '1'.

If it is big endian, it will be:

    +----+----+----+----+
    |0x00|0x00|0x00|0x01|
    +----+----+----+----+
    A
    |
   &x

so this one will be '0'.

*/
