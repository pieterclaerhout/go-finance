package finance_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pieterclaerhout/go-finance"
)

func TestExchangeRatesValid(t *testing.T) {

	rates, err := finance.ExchangeRates()

	assert.NoErrorf(t, err, "err should be nil, is: %v", err)
	assert.NotNilf(t, rates, "rates should not be nil")
	assert.NotEmpty(t, rates)

}

func TestExchangeRatesInvalidURL(t *testing.T) {

	finance.RatesURL = "ht&@-tp://:aa"
	defer resetRatesURL()

	rates, err := finance.ExchangeRates()

	assert.Error(t, err)
	assert.Empty(t, rates)

}

func TestExchangeRatesTimeout(t *testing.T) {

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

func TestExchangeRatesReadBodyError(t *testing.T) {

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

func TestExchangeRatesInvalidXML(t *testing.T) {

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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

func TestConvertRate(t *testing.T) {

	type test struct {
		name         string
		value        float64
		from         string
		to           string
		expectsError bool
	}

	var tests = []test{
		{"invalid", 1, "", "", true},
		{"invalid-from", 1, "EUR", "", true},
		{"invalid-to", 1, "", "EUR", true},
		{"valid-eur-eur", 2, "EUR", "EUR", false},
		{"valid-eur-usd", 2, "EUR", "USD", false},
		{"valid-aud-usd", 2, "AUD", "USD", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, err := finance.ConvertRate(tc.value, tc.from, tc.to)

			if tc.expectsError {
				assert.Zero(t, actual, "actual")
				assert.Error(t, err, "error")
			} else {
				assert.NotZero(t, actual, "actual")
				assert.NoError(t, err, "error")
				if tc.from != tc.to {
					assert.NotEqual(t, tc.value, actual, "should-not-be-equal")
				}
			}

		})
	}

}

func TestConvertRateInvalidURL(t *testing.T) {

	finance.RatesURL = "ht&@-tp://:aa"
	defer resetRatesURL()

	actual, err := finance.ConvertRate(1, "EUR", "USD")

	assert.Zero(t, actual, "actual")
	assert.Error(t, err, "error")

}

func resetRatesURL() {
	finance.RatesURL = finance.DefaultRatesURL
}
