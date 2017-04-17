package main

import (
	"github.com/flyingyizi/go-wiringPi/board"
	"github.com/flyingyizi/go-wiringPi/i2c"
)

func main() {
	info := board.GetBoardInfo()
	device := info.I2CDeviceName()

	d, err := i2c.Open(device)
	if err != nil {
		panic(err)
	}

	d, err = d.SetAddr(0xc0)
	// or opens a 10-bit address
	//d, err = d.SetAddr(i2c.TenBit(0x78))

	if err != nil {
		panic(err)
	}

	d.Close()

	_ = d
}
