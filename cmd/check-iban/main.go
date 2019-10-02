package main

import (
	"fmt"
	"os"

	"github.com/pieterclaerhout/go-finance"
)

func main() {

	info, err := finance.CheckIBAN("738120256174")
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}

	fmt.Println(info)

}
