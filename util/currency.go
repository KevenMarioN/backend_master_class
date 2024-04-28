package util

// Constants for all supported currencies
const (
	USD = "USD"
	BRL = "BRL"
	CAD = "CAD"
	EUR = "EUR"
)

// ISSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, BRL:
		return true
	}
	return false
}
