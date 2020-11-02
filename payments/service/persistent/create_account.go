package persistent

import (
	"context"
	"github.com/lightsgoout/fintech-go/payments/entity"
	"github.com/lightsgoout/fintech-go/payments/service"
	"github.com/lightsgoout/fintech-go/pkg/money"
	"strings"
)

func (s PaymentsService) CreateAccount(ctx context.Context, id entity.AccountID, balance money.Numeric, cur money.Currency) error {
	if balance.LessThan(money.NewNumericFromInt64(0)) {
		return service.ErrInsufficientFunds
	}

	if !money.IsKnownCurrency(cur) {
		return service.ErrIncompatibleCurrency
	}

	if id == "" {
		return service.ErrBadAccountID
	}

	const sql = `INSERT INTO account (id, currency, balance) VALUES (?id, ?currency, ?balance)`
	_, err := s.pg.ExecContext(ctx, sql, struct {
		Id       string `sql:"id"`
		Balance  string `sql:"balance"`
		Currency string `sql:"currency"`
	}{
		Id:       string(id),
		Balance:  balance.String(),
		Currency: string(cur),
	})
	if err != nil {
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "account_pkey"`) {
			return service.ErrAccountAlreadyExists
		}
		return NewInternalErrorFromDBError(err)
	}
	return nil
}
