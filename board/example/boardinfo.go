package main

import "fmt"

//go build -gcflags "-N -l"  blink.go

// LED Pin - BCM_GPIO 17.
const LED int = 17

func main() {
	fmt.Println("Raspberry Pi blink")
	info, base, err := board.GetBoardInfo()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("modelName:", info.ModelName())

	// Output: Raspberry Pi blink
}
