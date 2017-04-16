package main

import "golang.org/flyingyizi/go-wiringPi/i2c"

func main() {

	
	d, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x39)
	if err != nil {
		panic(err)
	}

	// opens a 10-bit address
	d, err = i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, i2c.TenBit(0x78))
	if err != nil {
		panic(err)
	}

	_ = d
}
