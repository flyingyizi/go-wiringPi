

## win7用久了,莫名奇妙,很多com端口号 都是使用中...把常用的端口都占完了....
打开CMD命令行，输入regedit打开注册表，找到HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\COM Name Arbiter，然后把ComDB删除,扫描检测硬件改动即可，不行的话重启PC即可

ref: https://github.com/johnlauer/serial-port-json-server

## How do I make serial work on the Raspberry Pi3

ref: https://raspberrypi.stackexchange.com/questions/45570/how-do-i-make-serial-work-on-the-raspberry-pi3

    This answer is still correct, and explains in more detail the nature of the changes, but most users of current Raspbian should just run sudo raspi-config Select Interfacing Options / Serial then specify if you want a Serial console (probably no) then if you want the Serial Port hardware enabled (probably yes). Then use /dev/serial0 in any code which accesses the Serial Port.

The BCM2837 on the Raspberry Pi3 has 2 UARTs (as did its predecessors), however to support the Bluetooth functionality the fully featured PL011 UART was moved from the header pins to the Bluetooth chip and the mini UART made available on header pins 8 & 10.

This has a number of consequences for users of the serial interface.

The /dev/ttyAMA0 previously used to access the UART now connects to Bluetooth.
The miniUART is now available on /dev/ttyS0.
In the latest operating system software there is a /dev/serial0 which selects the appropriate device so you can replace /dev/ttyAMA0 with /dev/serial0 and use the same software on the Pi3 and earlier models.

Unfortunately there are a number of other consequences:-

The mini UART is a secondary low throughput UART  
  intended to be used as a console.
The mini Uart has the following features:

1. • 7 or 8 bit operation.
2. • 1 start and 1 stop bit.
3. • No parities.
4. • Break generation.
5. • 8 symbols deep FIFOs for receive and transmit.
6. • SW controlled RTS, SW readable CTS.
7. • Auto flow control with programmable FIFO level.
8. • 16550 like registers.
9. • Baudrate derived from system clock.

There is no support for parity and the throughput is limited, but the latter should not affect most uses.

There is one killer feature "Baudrate derived from system clock" which makes the miniUART useless as the this clock can change dynamically e.g. if the system goes into reduced power or in low power mode.

Modifying the /boot/config.txt removes this dependency by adding the following line at the end:-

core_freq=250

This fixes the problem and appears to have little impact. The SPI clock frequency and ARM Timer are also dependent on the system clock.

    For some bizarre reason the default for Pi3 using the latest 4.4.9 kernel is to DISABLE UART. To enable it you need to change enable_uart=1 in /boot/config.txt. (This also fixes the core_freq so this is no longer necessary.)

Finally if you don't use Bluetooth (or have undemanding uses) it is possible to swap the ports back in Device Tree. There is a pi3-miniuart-bt and pi3-disable-bt module which are described in /boot/overlays/README.
