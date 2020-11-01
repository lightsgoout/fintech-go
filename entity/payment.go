package entity

import (
	"github.com/lightsgoout/fintech-go/pkg/money"
	"time"
)

type PaymentID int64

// Payment describes a transfer of money between one Account and another
type Payment struct {
	// id is a unique ID of this payment
	Id PaymentID

	Value PaymentValue
}

type PaymentValue struct {
	// ts is timestamp of transaction in UTC
	Time time.Time

	// from is AccountID from which money was taken
	From AccountID

	// to is AccountID to which money was transferred
	To AccountID

	// amount is amount of money transferred
	Amount money.Numeric

	Currency money.Currency
}
