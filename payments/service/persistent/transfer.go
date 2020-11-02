package persistent

import (
	"context"
	"github.com/lightsgoout/fintech-go/payments/entity"
	"github.com/lightsgoout/fintech-go/payments/service"
	"github.com/lightsgoout/fintech-go/pkg/money"
	"github.com/lightsgoout/fintech-go/pkg/postgres"
	"strings"
	"time"
)

func (s PaymentsService) Transfer(ctx context.Context, from, to entity.AccountID, amount money.Numeric, cur money.Currency) (entity.PaymentID, error) {
	// Freeze time so it would be consistent across all possible operations
	ts := time.Now().UTC()

	var paymentId entity.PaymentID

	// Open a transaction and lock both accounts for update
	// NOTE: we're not gonna modify accounts' primary keys,
	// so FOR NO KEY UPDATE is sufficient here and improves concurrency.
	// Also it's important to lock rows in deterministic order to prevent deadlocks.
	// We'll always lock the FROM account first.
	var lockOrder = [2]entity.AccountID{
		from,
		to,
	}

	err := postgres.NestedRunInTransaction(ctx, s.pg, func(tx postgres.Database) error {
		accounts := make(map[entity.AccountID]entity.Account, 2)
		for _, id := range lockOrder {
			account, err := s.getAccountWithLock(ctx, tx, id)
			if err != nil {
				if strings.Contains(err.Error(), "no rows in result set") {
					return service.ErrAccountDoesNotExist
				}
				return NewInternalErrorFromDBError(err)
			}
			accounts[id] = account
		}

		if accounts[from].Currency != accounts[to].Currency {
			return service.ErrIncompatibleCurrency
		}

		if accounts[from].Currency != cur {
			return service.ErrIncompatibleCurrency
		}

		newBalanceFrom := accounts[from].Balance.Sub(amount)
		newBalanceTo := accounts[to].Balance.Add(amount)
		if newBalanceFrom.LessThan(money.NewNumericFromInt64(0)) {
			return service.ErrInsufficientFunds
		}

		// With both accounts' locks acquired we can proceed to transfer the money.
		err := s.updateBalance(ctx, tx, from, newBalanceFrom)
		if err != nil {
			return NewInternalErrorFromDBError(err)
		}
		err = s.updateBalance(ctx, tx, to, newBalanceTo)
		if err != nil {
			return NewInternalErrorFromDBError(err)
		}

		// Create new Payment
		paymentId, err = s.createPayment(ctx, tx, entity.PaymentValue{
			Time:     ts,
			From:     from,
			To:       to,
			Amount:   amount,
			Currency: cur,
		})
		if err != nil {
			return NewInternalErrorFromDBError(err)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return paymentId, nil
}

func (s PaymentsService) updateBalance(ctx context.Context, tx postgres.Database, id entity.AccountID, newBalance money.Numeric) error {
	_, err := tx.ExecContext(ctx, `UPDATE account SET balance = ? WHERE id = ?`, newBalance.String(), id)
	return err
}

func (s PaymentsService) createPayment(ctx context.Context, tx postgres.Database, value entity.PaymentValue) (entity.PaymentID, error) {
	var result struct {
		Id int64 `sql:"id"`
	}
	const sql = `--payments_insert
		INSERT INTO payment
			(time, from_account_id, to_account_id, amount, currency)
		VALUES
			(?time, ?from_account_id, ?to_account_id, ?amount, ?currency)
		RETURNING
			id as id;
	`
	_, err := tx.QueryOneContext(ctx, &result, sql, struct {
		Time          time.Time `sql:"time"`
		FromAccountId string    `sql:"from_account_id"`
		ToAccountId   string    `sql:"to_account_id"`
		Amount        string    `sql:"amount"`
		Currency      string    `sql:"currency"`
	}{
		Time:          value.Time,
		FromAccountId: string(value.From),
		ToAccountId:   string(value.To),
		Amount:        value.Amount.String(),
		Currency:      string(value.Currency),
	})
	if err != nil {
		return 0, err
	}
	return entity.PaymentID(result.Id), nil
}

func (s PaymentsService) getAccountWithLock(ctx context.Context, tx postgres.Database, id entity.AccountID) (entity.Account, error) {
	var model struct {
		Id       string `sql:"id"`
		Currency string `sql:"currency"`
		Balance  string `sql:"balance"`
	}

	const sql = `SELECT id, currency, balance::text as balance FROM account WHERE id = ? FOR NO KEY UPDATE`

	_, err := tx.QueryOneContext(ctx, &model, sql, id)
	if err != nil {
		return entity.Account{}, err
	}

	return entity.Account{
		Id:       entity.AccountID(model.Id),
		Currency: money.Currency(model.Currency),
		Balance:  money.NewNumericFromStringMust(model.Balance),
	}, nil
}
