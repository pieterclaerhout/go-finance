package finance_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pieterclaerhout/go-finance"
)

const nameApple = "SRL APPLE RETAIL BELGIUM"
const addrApple = "Avenue du Port 86C/204\n1000 Bruxelles"

func TestCheck(t *testing.T) {

	type test struct {
		name                string
		vatNumber           string
		expectedCountryCode string
		expectedVATnumber   string
		expectedName        string
		expectedAddress     string
		expectedIsValid     bool
		expectedError       error
	}

	var tests = []test{
		{"empty", "", "", "", "", "", false, finance.ErrVATNumberTooShort},
		{"valid-spaces", "BE 0836 157 420", "BE", "0836157420", nameApple, addrApple, true, nil},
		{"valid-nospaces", "BE0836157420", "BE", "0836157420", nameApple, addrApple, true, nil},
		{"valid-dots", "BE 0836.157.420", "BE", "0836157420", nameApple, addrApple, true, nil},
		{"valid-ie", "IE6388047V", "IE", "6388047V", "GOOGLE IRELAND LIMITED", "3RD FLOOR, GORDON HOUSE, BARROW STREET, DUBLIN 4", true, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			time.Sleep(1 * time.Second)

			result, err := finance.CheckVAT(tc.vatNumber)

			if tc.expectedError != nil {

				assert.Nil(t, result)
				assert.Error(t, err)

			} else {

				assert.NotNil(t, result)
				assert.NoError(t, err)

				if result != nil {
					assert.Equal(t, tc.expectedCountryCode, result.CountryCode, "country-code")
					assert.Equal(t, tc.expectedVATnumber, result.VATNumber, "vat-number")
					assert.Equal(t, tc.expectedIsValid, result.IsValid, "is-valid")
					assert.Equal(t, tc.expectedName, result.Name, "name")
					assert.Equal(t, tc.expectedAddress, result.Address, "address")
				}

			}

		})
	}

}

func TestCheckVATInvalidURL(t *testing.T) {

	finance.VATServiceURL = "ht&@-tp://:aa"
	defer func() {
		finance.VATServiceURL = finance.DefaultVATServiceURL
	}()

	result, err := finance.CheckVAT("BE0836157420")

	assert.Nil(t, result)
	assert.Error(t, err)

}

func TestCheckVATTimeout(t *testing.T) {

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(500 * time.Millisecond)
			w.Header().Set("Content-Type", "text/xml")
			w.Write([]byte("hello"))
		}),
	)
	defer s.Close()

	finance.VATTimeout = 250 * time.Millisecond
	finance.VATServiceURL = s.URL
	defer func() {
		finance.VATTimeout = finance.DefaultVATTimeout
		finance.VATServiceURL = finance.DefaultVATServiceURL
	}()

	result, err := finance.CheckVAT("BE0836157420")

	assert.Nil(t, result)
	assert.Error(t, err)

}

func TestCheckVATReadBodyError(t *testing.T) {

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Length", "1")
		}),
	)
	defer s.Close()

	finance.VATServiceURL = s.URL
	defer func() {
		finance.VATServiceURL = finance.DefaultVATServiceURL
	}()

	result, err := finance.CheckVAT("BE0836157420")

	assert.Nil(t, result)
	assert.Error(t, err)

}

func TestCheckVATInvalidInput(t *testing.T) {

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte("INVALIDINPUT"))
		}),
	)
	defer s.Close()

	finance.VATServiceURL = s.URL
	defer func() {
		finance.VATServiceURL = finance.DefaultVATServiceURL
	}()

	result, err := finance.CheckVAT("BE0836157420")

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, finance.ErrVATnumberNotValid, err)

}

func TestCheckVATInvalidXML(t *testing.T) {

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte("<vies>"))
		}),
	)
	defer s.Close()

	finance.VATServiceURL = s.URL
	defer func() {
		finance.VATServiceURL = finance.DefaultVATServiceURL
	}()

	result, err := finance.CheckVAT("BE0836157420")

	assert.Nil(t, result)
	assert.Error(t, err)

}

func TestCheckVATSoapFault(t *testing.T) {

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><soap:Fault><faultcode>soap:Server</faultcode><faultstring>error</faultstring></soap:Fault></soap:Body></soap:Envelope>`))
		}),
	)
	defer s.Close()

	finance.VATServiceURL = s.URL
	defer func() {
		finance.VATServiceURL = finance.DefaultVATServiceURL
	}()

	result, err := finance.CheckVAT("BE0836157420")

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, finance.ErrVATserviceError+"error", err.Error(), "error-message")

}
