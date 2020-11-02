package money

import "strings"

type Currency string

func NewCurrency(raw string) Currency {
	return Currency(strings.ToUpper(raw))
}

func IsKnownCurrency(cur Currency) bool {
	switch cur {
	case "USD":
		return true
	case "EUR":
		return true
	case "RUB":
		return true
	}
	return false
}
