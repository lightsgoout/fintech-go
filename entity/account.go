package entity

import (
	"github.com/lightsgoout/fintech-go/pkg/money"
)

type AccountID string

type Account struct {
	Id       AccountID
	Balance  money.Numeric
	Currency money.Currency
}
