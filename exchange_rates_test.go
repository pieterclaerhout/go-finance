package ydfinance_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pieterclaerhout/go-ydfinance"
)

func Test_ExchangeRates_Valid(t *testing.T) {

	rates, err := ydfinance.ExchangeRates()

	assert.NoErrorf(t, err, "err should be nil, is: %v", err)
	assert.NotNilf(t, rates, "rates should not be nil")
	assert.NotEmpty(t, rates)

}

func Test_ExchangeRates_InvalidURL(t *testing.T) {

	ydfinance.RatesURL = "ht&@-tp://:aa"
	defer resetRatesURL()

	rates, err := ydfinance.ExchangeRates()

	assert.Error(t, err)
	assert.Empty(t, rates)

}

func Test_ExchangeRates_Timeout(t *testing.T) {

	ydfinance.DefaultTimeout = 250 * time.Millisecond

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(500 * time.Millisecond)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("hello"))
		}),
	)
	defer s.Close()

	ydfinance.RatesURL = s.URL
	defer resetRatesURL()

	rates, err := ydfinance.ExchangeRates()

	assert.Error(t, err)
	assert.Empty(t, rates)

}

func Test_ExchangeRates_ReadBodyError(t *testing.T) {

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "1")
		}),
	)
	defer s.Close()

	ydfinance.RatesURL = s.URL
	defer resetRatesURL()

	rates, err := ydfinance.ExchangeRates()

	assert.Error(t, err)
	assert.Empty(t, rates)

}

func Test_ExchangeRates_InvalidXML(t *testing.T) {

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("hello"))
		}),
	)
	defer s.Close()

	ydfinance.RatesURL = s.URL
	defer resetRatesURL()

	rates, err := ydfinance.ExchangeRates()

	assert.Error(t, err)
	assert.Empty(t, rates)

}

func resetRatesURL() {
	ydfinance.RatesURL = ydfinance.DefaultRatesURL
}
