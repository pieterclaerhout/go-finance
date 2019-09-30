package ydfinance

// exchangeRate defines the exchange rate
type exchangeRate struct {
	Time  string             `xml:"time"`
	Cubes []exchangeRateCube `xml:"Cube"`
}

// exchangeRateCube defines an exchange rate cube
type exchangeRateCube struct {
	TimedCubes []exchangeRateTimedCube `xml:"Cube"`
}

// exchangeRateTimedCube defines an exchange rate timed cube
type exchangeRateTimedCube struct {
	Time  string                     `xml:"time,attr"`
	Rates []exchangeRateCurrencyCube `xml:"Cube"`
}

// exchangeRateCurrencyCube defines an exchange rate currency cube
type exchangeRateCurrencyCube struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
}
