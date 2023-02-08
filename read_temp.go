package main

import (
	"fmt"
)

func readTemp() (float64, error) {

}

func main() {
	temp, err := readTemp()
	if err != nil {
		fmt.Println("Error cant read temp: ", err)
		return
	}
	fmt.Printf("Temp: %.2fC\n", temp)
}
