run-check-iban:
	go build -o check-iban github.com/pieterclaerhout/go-finance/cmd/check-iban
	./check-iban

run-check-vat:
	go build -o check-vat github.com/pieterclaerhout/go-finance/cmd/check-vat
	./check-vat
	
run-exchange-rates:
	go build -o exchange-rates github.com/pieterclaerhout/go-finance/cmd/exchange-rates
	./exchange-rates
