package service

import (
	"context"
	"github.com/lightsgoout/fintech-go/entity"
	"github.com/lightsgoout/fintech-go/pkg/money"
)

// PaymentsService is an interface containing all possible business operations for payments.
type PaymentsService interface {
	// Transfer sends money from one Account to another, atomically.
	Transfer(ctx context.Context, from, to entity.AccountID, amount money.Numeric, cur money.Currency) (entity.PaymentID, error)

	// GetPayments returns a list of transactions for a given AccountID in descending order (recent payments first).
	GetPayments(ctx context.Context, a entity.AccountID) ([]entity.Payment, error)

	// GetAccountList returns a list of possible AccountID's to trade with (matching the given Currency).
	GetAccountList(ctx context.Context, cur money.Currency) (entity.AccountID, error)
}
