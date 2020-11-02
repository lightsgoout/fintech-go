package service

import (
	"context"
	"github.com/lightsgoout/fintech-go/payments/entity"
	"github.com/lightsgoout/fintech-go/pkg/money"
)

// PaymentsService is an interface containing all possible business operations for payments.
type PaymentsService interface {
	// CreateAccount created new entity.Account
	CreateAccount(ctx context.Context, id entity.AccountID, balance money.Numeric, cur money.Currency) error

	// Transfer sends money from one entity.Account to another, atomically.
	Transfer(ctx context.Context, from, to entity.AccountID, amount money.Numeric, cur money.Currency) (entity.PaymentID, error)

	// GetPayments returns a list of transactions for a given AccountID in descending order (recent payments first).
	GetPayments(ctx context.Context, accountId entity.AccountID) ([]entity.Payment, error)

	// GetAccounts returns a list of possible AccountID's to trade with (matching the given Currency), ascending order.
	GetAccounts(ctx context.Context, cur money.Currency) ([]entity.AccountID, error)
}
