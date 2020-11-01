package money

import (
	"fmt"
	"golang.org/x/text/currency"
)

type Currency struct {
	unit currency.Unit
}

// NewCurrency takes raw currency string, validates it and returns Currency structure with a clean code
// e.g eur -> EUR, 1234 -> Err
func NewCurrency(s string) (Currency, error) {
	unit, err := currency.ParseISO(s)
	if err != nil {
		return Currency{}, fmt.Errorf("incorrect currency: %w", err)
	}

	return Currency{
		unit: unit,
	}, nil
}

// String returns currency ISO code, e.g EUR
func (c Currency) String() string {
	return c.unit.String()
}
