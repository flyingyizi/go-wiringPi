# Project Title

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.


### I2C Read/Write

The I2C read/write operation takes place as follows:

    1.The I2C master initiates the communication by sending a START condition followed by a 7-bit slave address and the 8th bit to indicate write (0)/ read (1)).

    2.The master releases the SDA and waits for an ACK from the slave device.

    3.If the slave exists on the bus, it responds with an ACK.

    4.The master continues in either transmit or receive mode (according to the read or write bit it sent), and the slave continues in its complementary mode (receive or transmit, respectively).

    5.The master terminates the data transmission by sending a STOP condition.

The following image shows a single byte read and write on an I2C slave device.

BYTE Read:
|Master | Start |  I2C slave address+read bit |       |       | NACK  | Stop |
|:-----:|:-----:|:---------------------------:|:-----:|:-----:|:-----:|:----:|
|slave  |       |                             |   ACK | Data  |       |      |
|       |       |                             |       |       |       |      |


BYTE Write:
|Master | Start |  I2C slave address+Writebit |       | Data  |       | Stop |
|:-----:|:-----:|:---------------------------:|:-----:|:-----:|:-----:|:----:|
|slave  |       |                             |   ACK |       |  ACK  |      |
|       |       |                             |       |       |       |      |



### I2C Register Read/Write

The I2C register read/write operation takes place as follows:

    1.The I2C master initiates the communication by sending a START condition followed by a 7-bit slave address and the 8th bit to indicate write (0)/ read (1).

    2.The master releases the SDA and waits for an ACK from slave device.

    3.If the slave exists on the bus, it responds with an ACK.

    4.Then, the master writes the register address of the slave it wants to access.

    5.Once the slave acknowledges the register address, the master sends the data byte with an ACK after each byte for write/read.

    The master terminates the data transmission by sending a STOP condition.

The following image shows a single byte read and write on a register present in the I2C slave device.




### IOCTL I2C_RDWR

This method allows for simultaneous read/write and sending an uninterrupted sequence of message. Not all i2c devices support this method.

Before performing i/o with this method, you should check whether the device supports this method using an ioctl I2C_FUNCS operation.

Using this method, you do not need to perform an ioctl I2C_SLAVE operation -- it is done behind the scenes using the information embedded in the messages.

### IOCTL SMBUS

This method of i/o is more powerful but the resulting code is more verbose. This method can be used if the device does not support the I2C_RDWR method.

Using this method, you do need to perform an ioctl I2C_SLAVE operation (or, if the device is busy, an I2C_SLAVE_FORCE operation).

### SYSFS I/O

This method uses the basic file i/o system calls read() and write(). Uninterrupted sequential operations are not possible using this method. This method can be used if the device does not support the I2C_RDWR method.

Using this method, you do need to perform an ioctl I2C_SLAVE operation (or, if the device is busy, an I2C_SLAVE_FORCE operation).

I can't think of any situation when this method would be preferable to others, unless you need the chip to be treated like a file.


### Prerequisites

What things you need to install the software and how to install them

```
Give examples
```

### Installing

A step by step series of examples that tell you have to get a development env running

Say what the step will be

```
Give the example
```

And repeat

```
until finished
```

End with an example of getting some data out of the system or using it for a little demo

## Running the tests

Explain how to run the automated tests for this system

### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

### And coding style tests

Explain what these tests test and why

go.exe test -timeout 30s -tags  -run ^Test_getEndian$



```
Give an example
```

## Deployment

Add additional notes about how to deploy this on a live system

## Built With

* [Dropwizard](http://www.dropwizard.io/1.0.2/docs/) - The web framework used
* [Maven](https://maven.apache.org/) - Dependency Management
* [ROME](https://rometools.github.io/rome/) - Used to generate RSS Feeds

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **Billie Thompson** - *Initial work* - [PurpleBooth](https://github.com/PurpleBooth)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone who's code was used
* Inspiration
* etc

* c sampleFull IOCTL Example

from : http://stackoverflow.com/questions/9974592/i2c-slave-ioctl-purpose
I haven't tested this example, but it shows the conceptual flow of writing to an i2c device.-- automatically detecting whether to use the ioctl I2C_RDWR or smbus technique.
```c
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include <errno.h>
#include <string.h>

#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>

#include <linux/i2c.h>
#include <linux/i2c-dev.h>
#include <sys/ioctl.h>

#define I2C_ADAPTER "/dev/i2c-0"
#define I2C_DEVICE  0x00

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

int i2c_ioctl_smbus_write (int fd, uint8_t dev, uint8_t regaddr, uint16_t *data)
{
    int i, j = 0;
    int ret;
    uint8_t *buf;

    buf = malloc(2 * (sizeof(data) / sizeof(data[0])));
    if (buf == NULL) {
        return -ENOMEM;
    }

    for (i = 0; i < (sizeof(data) / sizeof(data[0])); i ++) {
        buf[j ++] = (data[i] & 0xff00) >> 8;
        buf[j ++] = data[i] & 0xff;
    }

    struct i2c_smbus_ioctl_data payload = {
        .read_write = I2C_SMBUS_WRITE,
        .size = I2C_SMBUS_WORD_DATA,
        .command = regaddr,
        .data = (void *) buf,
    };

    ret = ioctl (fd, I2C_SLAVE_FORCE, dev);
    if (ret < 0)
    {
        ret = -errno;
        goto exit;
    }

    ret = ioctl (fd, I2C_SMBUS, &payload);
    if (ret < 0)
    {
        ret = -errno;
        goto exit;
    }

exit:
    free(buf);
    return ret;
}

int i2c_write (int fd, uint8_t dev, uint8_t regaddr, uint16_t *data)
{
    uint64_t funcs;

    if (ioctl(fd, I2C_FUNCS, &funcs) < 0) {
        return -errno;
    }

    if (funcs & I2C_FUNC_I2C) {
        return i2c_ioctl_write (fd, dev, regaddr, data);
    } else if (funcs & I2C_FUNC_SMBUS_WORD_DATA) {
        return i2c_ioctl_smbus_write (fd, dev, regaddr, data);
    } else {
        return -ENOSYS;
    }
}

int parse_args (uint8_t *regaddr, uint16_t *data, char *argv[])
{
    char *endptr;
    int i;

    *regaddr = (uint8_t) strtol(argv[1], &endptr, 0);
    if (errno || endptr == argv[1]) {
        return -1;
    }

    for (i = 0; i < sizeof(data) / sizeof(data[0]); i ++) {
        data[i] = (uint16_t) strtol(argv[i + 2], &endptr, 0);
        if (errno || endptr == argv[i + 2]) {
            return -1;
        }
    }

    return 0;
}

void usage (int argc, char *argv[])
{
    fprintf(stderr, "Usage: %s regaddr data [data]*\n", argv[0]);
    fprintf(stderr, "  regaddr   The 8-bit register address to write to.\n");
    fprintf(stderr, "  data      The 16-bit data to be written.\n");
    exit(-1);
}

int main (int argc, char *argv[])
{
    uint8_t regaddr;
    uint16_t *data;
    int fd;
    int ret = 0;

    if (argc < 3) {
        usage(argc, argv);
    }

    data = malloc((argc - 2) * sizeof(uint16_t));
    if (data == NULL) {
        fprintf (stderr, "%s.\n", strerror(ENOMEM));
        return -ENOMEM;
    }

    if (parse_args(&regaddr, data, argv) != 0) {
        free(data);
        usage(argc, argv);
    }

    fd = open(I2C_ADAPTER, O_RDWR | O_NONBLOCK);
    ret = i2c_write(fd, I2C_DEVICE, regaddr, data);
    close(fd);

    if (ret) {
        fprintf (stderr, "%s.\n", strerror(-ret));
    }

    free(data);

    return ret;
}
```