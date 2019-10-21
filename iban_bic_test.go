package finance_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pieterclaerhout/go-finance"
)

func TestCheckIBAN(t *testing.T) {

	type test struct {
		number           string
		expectedBankName string
		expectedIBAN     string
		expectedBIC      string
		expectsError     bool
	}

	var tests = []test{
		{"738-1202561-74", "KBC Bank", "BE16 7381 2025 6174", "KRED BE BB", false},
		{"738120256174", "KBC Bank", "BE16 7381 2025 6174", "KRED BE BB", false},
		{"7381202561-74", "KBC Bank", "BE16 7381 2025 6174", "KRED BE BB", false},
		{"738-AAAAAAA-74", "", "", "", true},
		{"738-AAAAAAA-AA", "", "", "", true},
		{"", "", "", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.number, func(t *testing.T) {

			info, err := finance.CheckIBAN(tc.number)

			if tc.expectsError {

				assert.Error(t, err, "error")
				assert.Nil(t, info, "info")

			} else {

				assert.NoError(t, err, "error")
				assert.NotNil(t, info, "info")

				if info != nil {
					assert.Equal(t, tc.expectedBankName, info.BankName, "bank-name")
					assert.Equal(t, tc.expectedIBAN, info.IBAN, "IBAN")
					assert.Equal(t, tc.expectedBIC, info.BIC, "BIC")
				}

			}

		})
	}

}

func TestCheckIBANInvalidURL(t *testing.T) {

	finance.IBANBICServiceURL = "ht&@-tp://:aa"
	defer func() {
		finance.IBANBICServiceURL = finance.DefaultIBANBICServiceURL
	}()

	result, err := finance.CheckIBAN("738120256174")

	assert.Nil(t, result, "result")
	assert.Error(t, err, "error")

}

func TestCheckIBANTimeout(t *testing.T) {

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(500 * time.Millisecond)
			w.Header().Set("Content-Type", "text/xml")
			w.Write([]byte("hello"))
		}),
	)
	defer s.Close()

	finance.IBANBICTimeout = 250 * time.Millisecond
	finance.IBANBICServiceURL = s.URL
	defer func() {
		finance.IBANBICTimeout = finance.DefaultIBANBICTimeout
		finance.IBANBICServiceURL = finance.DefaultIBANBICServiceURL
	}()

	result, err := finance.CheckIBAN("738120256174")

	assert.Nil(t, result, "result")
	assert.Error(t, err, "error")

}

func TestCheckIBANReadBodyError(t *testing.T) {

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(`<string xmlns="http://tempuri.org/">KBC Bank</string>`))
		}),
	)
	defer s.Close()

	finance.IBANBICServiceURL = s.URL
	defer func() {
		finance.IBANBICServiceURL = finance.DefaultIBANBICServiceURL
	}()

	result, err := finance.CheckIBAN("738120256174")

	assert.Nil(t, result, "result")
	assert.Error(t, err, "error")

}

func TestCheckIBANInvalidBody(t *testing.T) {

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Length", "1")
		}),
	)
	defer s.Close()

	finance.IBANBICServiceURL = s.URL
	defer func() {
		finance.IBANBICServiceURL = finance.DefaultIBANBICServiceURL
	}()

	result, err := finance.CheckIBAN("738120256174")

	assert.Nil(t, result, "result")
	assert.Error(t, err, "error")

}

func TestCheckIBANPartialFail(t *testing.T) {

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.RequestURI, "BBANtoIBANandBIC") {
				w.Header().Set("Content-Length", "1")
				return
			}
			w.Write([]byte(`<string xmlns="http://tempuri.org/">KBC Bank</string>`))
		}),
	)
	defer s.Close()

	finance.IBANBICServiceURL = s.URL
	defer func() {
		finance.IBANBICServiceURL = finance.DefaultIBANBICServiceURL
	}()

	result, err := finance.CheckIBAN("738120256174")

	assert.Nil(t, result, "result")
	assert.Error(t, err, "error")

}
