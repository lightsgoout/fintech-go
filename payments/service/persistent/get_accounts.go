package persistent

import (
	"context"
	"github.com/lightsgoout/fintech-go/payments/entity"
	"github.com/lightsgoout/fintech-go/payments/service"
	"github.com/lightsgoout/fintech-go/pkg/money"
)

func (s PaymentsService) GetAccounts(ctx context.Context, cur money.Currency) ([]entity.AccountID, error) {
	if !money.IsKnownCurrency(cur) {
		return nil, service.ErrIncompatibleCurrency
	}
	result, err := s.getAccounts(ctx, cur)
	if err != nil {
		return nil, NewInternalErrorFromDBError(err)
	}
	return result, nil
}

func (s PaymentsService) getAccounts(ctx context.Context, cur money.Currency) ([]entity.AccountID, error) {
	const sql = `SELECT id FROM account WHERE currency = ? ORDER BY id ASC`
	var rows []entity.AccountID
	_, err := s.pg.QueryContext(ctx, &rows, sql, cur)
	return rows, err
}
