package persistent

import (
	"context"
	"github.com/lightsgoout/fintech-go/payments/entity"
	"github.com/lightsgoout/fintech-go/payments/service"
	"github.com/lightsgoout/fintech-go/pkg/money"
	"time"
)

func (s PaymentsService) GetPayments(ctx context.Context, accountId entity.AccountID) ([]entity.Payment, error) {
	// Check account exists
	exists, err := s.accountExists(ctx, accountId)
	if err != nil {
		return nil, NewInternalErrorFromDBError(err)
	}
	if !exists {
		return nil, service.ErrAccountDoesNotExist
	}
	result, err := s.getPayments(ctx, accountId)
	if err != nil {
		return nil, NewInternalErrorFromDBError(err)
	}
	return result, nil
}

func (s PaymentsService) accountExists(ctx context.Context, accountId entity.AccountID) (bool, error) {
	const sql = `SELECT exists(select id from account where id = ?) as exists`
	var result struct {
		Exists bool `sql:"exists"`
	}
	_, err := s.pg.QueryOneContext(ctx, &result, sql, accountId)
	return result.Exists, err
}

func (s PaymentsService) getPayments(ctx context.Context, accountId entity.AccountID) ([]entity.Payment, error) {
	const sql = `--payments_get
		SELECT * FROM (
			SELECT 
				id, 
				time, 
				from_account_id, 
				to_account_id, 
				amount::text as amount, 
				currency,
				true as outgoing
			FROM payment WHERE from_account_id = ?
			UNION
			SELECT 
				id, 
				time, 
				from_account_id, 
				to_account_id, 
				amount::text as amount, 
				currency,
				false as outgoing
			FROM payment WHERE to_account_id = ?
		) x ORDER BY time DESC`

	var rows []struct {
		Id       int64     `sql:"id"`
		Time     time.Time `sql:"time"`
		From     string    `pg:"from_account_id"`
		To       string    `pg:"to_account_id"`
		Amount   string    `sql:"amount"`
		Currency string    `sql:"currency"`
		Outgoing bool      `sql:"outgoing"`
	}
	_, err := s.pg.QueryContext(ctx, &rows, sql, accountId, accountId)
	if err != nil {
		return nil, err
	}
	result := make([]entity.Payment, 0, len(rows))
	for _, r := range rows {
		result = append(result, entity.Payment{
			Id: entity.PaymentID(r.Id),
			Value: entity.PaymentValue{
				Time:     r.Time,
				From:     entity.AccountID(r.From),
				To:       entity.AccountID(r.To),
				Amount:   money.NewNumericFromStringMust(r.Amount),
				Currency: money.Currency(r.Currency),
				Outgoing: r.Outgoing,
			},
		})
	}
	return result, nil
}
