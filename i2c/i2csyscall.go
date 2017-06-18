//+build linux

//go:generate sh i2csyscall.sh $GOFILE $GOOS $GOARCH

package i2c
