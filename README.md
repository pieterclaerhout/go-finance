# go-finance

[![Go Report Card](https://goreportcard.com/badge/github.com/pieterclaerhout/go-finance)](https://goreportcard.com/report/github.com/pieterclaerhout/go-finance)
[![Documentation](https://godoc.org/github.com/pieterclaerhout/go-finance?status.svg)](http://godoc.org/github.com/pieterclaerhout/go-finance)
[![license](https://img.shields.io/badge/license-Apache%20v2-orange.svg)](https://github.com/pieterclaerhout/go-finance/raw/master/LICENSE)
[![GitHub version](https://badge.fury.io/gh/pieterclaerhout%2Fgo-finance.svg)](https://badge.fury.io/gh/pieterclaerhout%2Fgo-finance)
[![GitHub issues](https://img.shields.io/github/issues/pieterclaerhout/go-finance.svg)](https://github.com/pieterclaerhout/go-finance/issues)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fpieterclaerhout%2Fgo-finance.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fpieterclaerhout%2Fgo-finance?ref=badge_shield)

This is a [Golang](https://golang.org) library which contains finance related functions.

## Exchange Rates

The following example explains how to use this package to retrieve the exchange rates from [ECB](https://www.ecb.europa.eu):

```go
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
```

## Checking VAT Numbers


You can also VAT numbers via the [VIES service](http://ec.europa.eu/taxation_customs/vies/vatRequest.html). The following sample code shows how to do this:


```go
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
```

## IBAN & BIC

There is also a function which converts a regular Belgian Bank Account Number to it's IBAN / BIC equivalent:

```go
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
```

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fpieterclaerhout%2Fgo-finance.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fpieterclaerhout%2Fgo-finance?ref=badge_large)