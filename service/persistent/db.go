package persistent

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/lightsgoout/fintech-go/entity"
	"github.com/lightsgoout/fintech-go/pkg/money"
	"time"
)

// PaymentsService implements service.PaymentsService interface by using persistent storage.
type PaymentsService struct {
	pg *pg.DB
}

func (db PaymentsService) Transfer(ctx context.Context, from, to entity.AccountID, amount money.Numeric, cur money.Currency) (entity.PaymentID, error) {
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

	err := db.pg.RunInTransaction(ctx, func(tx *pg.Tx) error {
		accounts := make(map[entity.AccountID]entity.Account, 2)
		for _, id := range lockOrder {
			_, err := tx.QueryOne(accounts[id], `SELECT * FROM account WHERE id = ? FOR NO KEY UPDATE`, id)
			if err != nil {
				// TODO: account does not exist?
				return fmt.Errorf("database error: %w", err)
			}
		}

		if accounts[from].Currency != accounts[to].Currency {
			return errIncompatibleCurrency
		}

		if accounts[from].Currency != cur {
			return errIncompatibleCurrency
		}

		newBalanceFrom := accounts[from].Balance.Sub(amount)
		newBalanceTo := accounts[to].Balance.Add(amount)
		if newBalanceFrom.LessThan(money.NewNumericFromInt64(0)) {
			return errInsufficientFunds
		}

		// With both accounts' locks acquired we can proceed to transfer the money.
		err := db.updateBalance(tx, from, newBalanceFrom)
		if err != nil {
			return fmt.Errorf("database error: %w", err)
		}
		err = db.updateBalance(tx, to, newBalanceTo)
		if err != nil {
			return fmt.Errorf("database error: %w", err)
		}

		// Create new Payment
		paymentId, err = db.createPayment(tx, entity.PaymentValue{
			Time:     ts,
			From:     from,
			To:       to,
			Amount:   amount,
			Currency: cur,
		})
		if err != nil {
			return fmt.Errorf("database error: %w", err)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return paymentId, nil
}

func (db PaymentsService) updateBalance(tx *pg.Tx, id entity.AccountID, newBalance money.Numeric) error {
	_, err := tx.Exec(`UPDATE account SET balance = ? WHERE id = ?`, newBalance.String(), id)
	return err
}

func (db PaymentsService) createPayment(tx *pg.Tx, value entity.PaymentValue) (entity.PaymentID, error) {
	var result entity.PaymentID
	const sql = `--payments_insert
		INSERT INTO payment
			(time, from_account_id, to_account_id, amount, currency)
		VALUES
			(?time, ?from_account_id, ?to_account_id, ?amount, ?currency)
		RETURNING
			id;
	`
	_, err := tx.QueryOne(&result, sql, struct {
		time          time.Time `sql:"time"`
		fromAccountId string    `sql:"from_account_id"`
		toAccountId   string    `sql:"to_account_id"`
		amount        string    `sql:"amount"`
		currency      string    `sql:"currency"`
	}{
		time:          value.Time,
		fromAccountId: string(value.From),
		toAccountId:   string(value.To),
		amount:        value.Amount.String(),
		currency:      value.Currency.String(),
	})
	if err != nil {
		return 0, err
	}
	return result, nil
}
