package finance

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// DefaultRatesURL defines the default URL to fetch the exchange rates from
const DefaultRatesURL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

// RatesURL is the URL where to fetch the rates from
var RatesURL = DefaultRatesURL

// DefaultTimeout is the default tiemout for the HTTP client
var DefaultTimeout = 5 * time.Second

// ExchangeRates returs the list exchange rates
func ExchangeRates() (map[string]float64, error) {

	var rates exchangeRate

	ratesMap := make(map[string]float64, 0)

	client := &http.Client{}
	client.Timeout = DefaultTimeout

	resp, err := client.Get(RatesURL)
	if err != nil {
		return ratesMap, err
	}
	defer resp.Body.Close()

	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ratesMap, err
	}

	err = xml.Unmarshal(rawData, &rates)
	if err != nil {
		return ratesMap, err
	}

	ratesMap["EUR"] = 1

	for _, cube := range rates.Cubes {
		for _, timedCube := range cube.TimedCubes {
			for _, rate := range timedCube.Rates {
				ratesMap[strings.ToUpper(rate.Currency)] = rate.Rate
			}
		}
	}

	return ratesMap, nil

}
