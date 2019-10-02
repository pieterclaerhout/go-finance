package finance

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// DefaultIBANBICServiceURL is the default IBANBIC service URL to use
const DefaultIBANBICServiceURL = "https://www.ibanbic.be/IBANBIC.asmx"

// IBANBICServiceURL is the SOAP URL to be used when checking a bank account number
var IBANBICServiceURL = DefaultIBANBICServiceURL

// DefaultIBANBICTimeout is the default timeout to use when checking the bank account number
const DefaultIBANBICTimeout = 5 * time.Second

// IBANBICTimeout is the timeout to use when checking the bank account number
var IBANBICTimeout = DefaultIBANBICTimeout

var (
	// ErrIBANBICServiceUnreachable is the error returned when the IBANBIC service is unreachable
	ErrIBANBICServiceUnreachable = errors.New("IBANBIC service is unreachable")

	// ErrIBANBICInvalidInput is the error returned when the bank account number is invalid
	ErrIBANBICInvalidInput = errors.New("Number is not a valid bank account number")

	// ErrIBANBICServiceError is the error returned when we get a non-standard error from the IBANBIC service
	ErrIBANBICServiceError = "IBANBIC service returns an error: "
)

// IBANBICInfo contains the info about a Belgian Bank Account number
type IBANBICInfo struct {
	BBAN     string // The Belgian Bank Account Number
	BankName string // The name of the bank which issues the account
	IBAN     string // The IBAN number of the bank account
	BIC      string // The Bank Identification Code of the bank
}

// CheckIBAN checks the Bank Account Number and returns the IBAN and BIC information
func CheckIBAN(number string) (*IBANBICInfo, error) {

	if len(number) == 0 {
		return nil, ErrIBANBICInvalidInput
	}

	result := &IBANBICInfo{
		BBAN: number,
	}

	bankName, err := performIBANBICRequest("BBANtoBANKNAME", number)
	if err != nil {
		return nil, err
	}
	result.BankName = bankName

	ibanAndBic, err := performIBANBICRequest("BBANtoIBANandBIC", number)
	if err != nil {
		return nil, err
	}

	ibanBicParts := strings.SplitN(ibanAndBic, "#", 2)
	if len(ibanBicParts) < 2 {
		return nil, errors.New(ErrIBANBICServiceError + "Failed to get BIC and IBAN code")
	}

	result.IBAN = ibanBicParts[0]
	result.BIC = ibanBicParts[1]

	return result, nil

}

func performIBANBICRequest(action string, value string) (string, error) {

	url := IBANBICServiceURL + "/" + url.PathEscape(action) + "?Value=" + url.QueryEscape(value)

	client := http.Client{
		Timeout: VATTimeout,
	}

	res, err := client.Get(url)
	if err != nil {
		return "", ErrIBANBICServiceUnreachable
	}
	defer res.Body.Close()

	xmlRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	xmlString := string(xmlRes)
	if strings.Contains(xmlString, "Exception") {
		exceptionParts := strings.Split(xmlString, "\n")
		return "", errors.New(ErrIBANBICServiceError + strings.TrimSpace(exceptionParts[0]))
	}

	var result string
	if err := xml.Unmarshal(xmlRes, &result); err != nil {
		return "", err
	}

	return result, nil

}
