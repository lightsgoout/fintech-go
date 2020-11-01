package api

import (
	"github.com/lightsgoout/fintech-go/entity"
	"github.com/lightsgoout/fintech-go/pkg/money"
)

type transferRequest struct {
	from     entity.AccountID
	to       entity.AccountID
	amount   money.Numeric
	currency money.Currency
}
