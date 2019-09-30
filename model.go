package ydfinance

// ExchangeRate defines the exchange rate
type ExchangeRate struct {
	Time  string             `xml:"time"`
	Cubes []ExchangeRateCube `xml:"Cube"`
}

// ExchangeRateCube defines an exchange rate cube
type ExchangeRateCube struct {
	TimedCubes []ExchangeRateTimedCube `xml:"Cube"`
}

// ExchangeRateTimedCube defines an exchange rate timed cube
type ExchangeRateTimedCube struct {
	Time  string                     `xml:"time,attr"`
	Rates []ExchangeRateCurrencyCube `xml:"Cube"`
}

// ExchangeRateCurrencyCube defines an exchange rate currency cube
type ExchangeRateCurrencyCube struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
}
