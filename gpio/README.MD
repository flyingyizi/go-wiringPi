
# get baoard info

the code refer https://github.com/AndrewFromMelbourne/raspberry_pi_revision  and wiringPI archive

## Prior to Raspberry Pi 2

that holds the revision number of the Raspberry Pi.
Known revisions (prior to the Raspberry Pi 2) are:

    | Revision |  Model  | PCB Rev | Memory | Manufacture |
    |:--------:|:-------:|:-------:|:------:|:------------|
    |   0000   |         |         |        |             |
    |   0001   |         |         |        |             |
    |   0002   |    B    |    1    | 256 MB |  EGOMAN     |
    |   0003   |    B    |    1-1  | 256 MB |   EGOMAN    |
    |   0004   |    B    |    1-2  | 256 MB |   Sony      |
    |   0005   |    B    |    1-2  | 256 MB |   EGOMAN    |
    |   0006   |    B    |    1-2  | 256 MB |   EGOMAN    |
    |   0007   |    A    |    1-2  | 256 MB |   Egoman    |
    |   0008   |    A    |    1-2  | 256 MB |   Sony      |
    |   0009   |    A    |    1-2  | 256 MB |   Egoman    |
    |   000d   |    B    |    1-2  | 512 MB |   Egoman    |
    |   000e   |    B    |    1-2  | 512 MB |   Sony      |
    |   000f   |    B    |    1-2  | 512 MB |   Egoman    |
    |   0010   |    B+   |    1-2  | 512 MB |   Sony      |
    |   0011   | compute |    1-1  | 512 MB |   Sony      |
    |   0012   |    A+   |    1-1  | 256 MB |   Sony      |
    |   0013   |    B+   |    1-2  | 512 MB |   Embest    |
    |   0014   | compute |    1-1  | 512 MB |   Embest    |
    |   0015   |    A+   |    1-1  | 256 MB |   Embest    |
    |   0016   |    A+   |    1-1  | 256 MB |   Embest    |
    |   0017   | compute |    1-1  | 512 MB |   sony      |
    |   0018   | A+      |    1-1  | 256 MB |   sony      |
    |   0019   |    B+   |    1-2  | 512 MB |   egoman    |
    |   001a   | compute |    1-1  | 512 MB |   egoman    |
    |   001b   | A+      |    1-1  | 256 MB |   egoman    |

If the Raspberry Pi has been over-volted (voiding the warranty) the
revision number will have 100 at the front. e.g. 1000002.

## Raspberry Pi 2 / Model Zero

With the release of the Raspberry Pi 2, there is a new encoding of the
Revision field in /proc/cpuinfo. The bit fields are as follows


|``FEDC``|``BA98``|``7654``|``3210``|``FEDC``|``BA98``|``7654``|``3210``|
|---:|---:|---:|---:|---:|---:|---:|---:|
|    |  ``HG``|``FEEE``|``DDDD``|``CCCC``|``BBBB``|``BBBB``|``AAAA``|
|``1098``|``7654``|``3210``|``9876``|``5432``|``1098``|``7654``|``3210``|

| # | bits  |   contains   | values                                                                                            |
|:-:|:-----:|:-------------|:--------------------------------------------------------------------------------------------------|
| A | 00-03 | PCB Revision | (the pcb revision number)                                                                         |
| B | 04-11 | Model name   | A, B, A+, B+, B Pi2, Alpha, Compute Module, unknown, Pi3, Zero, Compute Module 3, unknown, Zero W |
| C | 12-15 | Processor    | BCM2835, BCM2836, BCM2837                                                                         |
| D | 16-19 | Manufacturer | Sony, Egoman, Embest, Sony Japan, Embest                                                          |
| E | 20-22 | Memory size  | 256 MB, 512 MB, 1024 MB                                                                           |
| F | 23-23 | encoded flag | (if set, revision is a bit field)                                                                 |
| G | 24-24 | waranty bit  | (if set, warranty void - Pre Pi2)                                                                 |
| H | 25-25 | waranty bit  | (if set, warranty void - Post Pi2)                                                                |

Also, due to some early issues the warranty bit has been move from bit
24 to bit 25 of the revision number (i.e. 0x2000000).


## 
Peripheral Offsets for the Raspberry Pi
|Peripheral  | Offset     | Address    | Description       |  C Offset Macro     |     
|------------|------------|------------|-------------------|---------------------|
|Base        | 0x00000000 | 0x20000000 | Starting address  |  BCM2708_PERI_BASE  |
|PADS_GPIO   | 0x00100000 | 0x20100000 | PADS base         |  PADS_GPIO_BASE     |
|GPIO 00..27 | 0x0010002C | 0x2010002C | GPIO 00..27 pads  |  PADS_GPIO_00_27    |
|GPIO 28..45 | 0x00100030 | 0x20100030 | GPIO 28..45 pads  |  PADS_GPIO_28_45    |
|GPIO 46..53 | 0x00100034 | 0x20100034 | GPIO 46..53 pads  |  PADS_GPIO_46_53    |
|Clock       | 0x00101000 | 0x20101000 | Clock registers   |  ClkBase            |
|GPIO        | 0x00200000 | 0x20200000 | GPIO registers    |  GpioBase           |
|GPPUD       | 0x00200025 | 0x20200025 | Pull-up enable    |                     |
|GPPUDCLK0   | 0x00200026 | 0x20200026 | Pull-up clock 0   |                     |
|GPPUDCLK1   | 0x00200027 | 0x20200027 | Pull-up clock 1   |                     |
|PWM         | 0x0020C000 | 0x2020C000 | PWM registers     |  PwmBase            |           
|            |            |            |                   |                     |