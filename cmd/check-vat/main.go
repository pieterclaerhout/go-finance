package main

import (
	"fmt"
	"os"

	"github.com/pieterclaerhout/go-finance"
)

func main() {

	info, err := finance.CheckVAT("BE0836157420")
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}

	fmt.Println(info)

}
