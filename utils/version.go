package utils

import (
	"fmt"
)

const (
	Version = "0.0.1"
)

func PrintVersion() {
	fmt.Println("dmux version", Version)
}
