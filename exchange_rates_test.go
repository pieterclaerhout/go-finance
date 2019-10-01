package finance_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pieterclaerhout/go-finance"
)

func Test_ExchangeRates_Valid(t *testing.T) {

	rates, err := finance.ExchangeRates()

	assert.NoErrorf(t, err, "err should be nil, is: %v", err)
	assert.NotNilf(t, rates, "rates should not be nil")
	assert.NotEmpty(t, rates)

}

func Test_ExchangeRates_InvalidURL(t *testing.T) {

	finance.RatesURL = "ht&@-tp://:aa"
	defer resetRatesURL()

	rates, err := finance.ExchangeRates()

	assert.Error(t, err)
	assert.Empty(t, rates)

}

func Test_ExchangeRates_Timeout(t *testing.T) {

	finance.DefaultTimeout = 250 * time.Millisecond

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(500 * time.Millisecond)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("hello"))
		}),
	)
	defer s.Close()

	finance.RatesURL = s.URL
	defer resetRatesURL()

	rates, err := finance.ExchangeRates()

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

	finance.RatesURL = s.URL
	defer resetRatesURL()

	rates, err := finance.ExchangeRates()

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

	finance.RatesURL = s.URL
	defer resetRatesURL()

	rates, err := finance.ExchangeRates()

	assert.Error(t, err)
	assert.Empty(t, rates)

}

func resetRatesURL() {
	finance.RatesURL = finance.DefaultRatesURL
}
