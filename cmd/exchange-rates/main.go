package main

import (
	"fmt"
	"os"

	"github.com/pieterclaerhout/go-finance"
)

func main() {

	rates, err := finance.ExchangeRates()
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}

	for currency, rate := range rates {
		fmt.Println(currency, "-> €1 =", rate)
	}

}
