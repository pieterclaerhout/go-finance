package ydfinance

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"
)

// DefaultRatesURL defines the default URL to fetch the exchange rates from
const DefaultRatesURL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

// RatesURL is the URL where to fetch the rates from
var RatesURL = DefaultRatesURL

// DefaultTimeout is the default tiemout for the HTTP client
var DefaultTimeout = 5 * time.Second

// ExchangeRates returs the list exchange rates
func ExchangeRates() (ExchangeRate, error) {

	var rates ExchangeRate

	client := &http.Client{}
	client.Timeout = DefaultTimeout

	resp, err := client.Get(RatesURL)
	if err != nil {
		return rates, err
	}
	defer resp.Body.Close()

	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rates, err
	}

	err = xml.Unmarshal(rawData, &rates)
	if err != nil {
		return rates, err
	}

	return rates, nil

}
