package finance

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// VATInfo is the info returned about a VAT number
type VATInfo struct {
	CountryCode string // The country code
	VATNumber   string // The VAT number
	IsValid     bool   // A boolean indicating if the VAT number is valid
	Name        string // The name linked to the VAT number
	Address     string // The address linked to the VAT number
}

// DefaultVATServiceURL is the default VAT service URL to use
const DefaultVATServiceURL = "http://ec.europa.eu/taxation_customs/vies/services/checkVatService"

// VATServiceURL is the SOAP URL to be used when checking a VAT number
var VATServiceURL = DefaultVATServiceURL

// DefaultVATTimeout is the default timeout to use when checking the VAT service
const DefaultVATTimeout = 5 * time.Second

// VATTimeout is the timeout to use when checking the VAT service
var VATTimeout = DefaultVATTimeout

var (
	// ErrVATnumberNotValid is the error returned when the VAT number is invalid
	ErrVATnumberNotValid = errors.New("VAT number is not valid")

	// ErrVATserviceUnreachable is the error returned when the VIES service is unreachable
	ErrVATserviceUnreachable = errors.New("VAT number validation service is unreachable")

	// ErrVATNumberTooShort is the error returned when the VAT number is too short
	ErrVATNumberTooShort = errors.New("VAT number is too short")

	// ErrVATserviceError is the error returned when we get a non-standard error from the VAT service
	ErrVATserviceError = "VAT number validation service returns an error: "
)

// CheckVAT checks the VAT number and returns the data
func CheckVAT(vatNumber string) (*VATInfo, error) {

	vatNumber = sanitizeVatNumber(vatNumber)

	e, err := buildViewEnvelope(vatNumber)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: VATTimeout,
	}

	eb := bytes.NewBufferString(e)
	res, err := client.Post(VATServiceURL, "text/xml;charset=UTF-8", eb)
	if err != nil {
		return nil, ErrVATserviceUnreachable
	}
	defer res.Body.Close()

	xmlRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if bytes.Contains(xmlRes, []byte("INVALID_INPUT")) {
		return nil, ErrVATnumberNotValid
	}

	// return nil, errors.New(string(xmlRes))

	var rd struct {
		XMLName xml.Name `xml:"Envelope"`
		Soap    struct {
			XMLName xml.Name `xml:"Body"`
			Soap    struct {
				XMLName     xml.Name `xml:"checkVatResponse"`
				CountryCode string   `xml:"countryCode"`
				VATnumber   string   `xml:"vatNumber"`
				Valid       bool     `xml:"valid"`
				Name        string   `xml:"name"`
				Address     string   `xml:"address"`
			}
			SoapFault struct {
				XMLName string `xml:"Fault"`
				Code    string `xml:"faultcode"`
				Message string `xml:"faultstring"`
			}
		}
	}
	if err := xml.Unmarshal(xmlRes, &rd); err != nil {
		return nil, err
	}

	if rd.Soap.SoapFault.Message != "" {
		return nil, errors.New(ErrVATserviceError + rd.Soap.SoapFault.Message)
	}

	r := &VATInfo{
		CountryCode: rd.Soap.Soap.CountryCode,
		VATNumber:   rd.Soap.Soap.VATnumber,
		IsValid:     rd.Soap.Soap.Valid,
	}

	if r.IsValid {
		r.Name = rd.Soap.Soap.Name
		r.Address = rd.Soap.Soap.Address
	}

	return r, nil

}

// sanitizeVatNumber removes all white space from a string
func sanitizeVatNumber(vatNumber string) string {
	vatNumber = strings.TrimSpace(vatNumber)
	vatNumber = strings.ReplaceAll(vatNumber, " ", "")
	vatNumber = strings.ReplaceAll(vatNumber, ".", "")
	return vatNumber
}

// buildViewEnvelope parses envelope template
func buildViewEnvelope(vatNumber string) (string, error) {

	if len(vatNumber) < 3 {
		return "", ErrVATNumberTooShort
	}

	envelope := `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:v1="http://schemas.conversesolutions.com/xsd/dmticta/v1">
<soapenv:Header/>
<soapenv:Body>
  <checkVat xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
    <countryCode>` + strings.ToUpper(vatNumber[0:2]) + `</countryCode>
    <vatNumber>` + strings.ToUpper(vatNumber[2:]) + `</vatNumber>
  </checkVat>
</soapenv:Body>
</soapenv:Envelope>`

	return envelope, nil

}
