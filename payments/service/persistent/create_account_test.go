package persistent

import (
	"errors"
	"github.com/lightsgoout/fintech-go/payments/service"
	"github.com/lightsgoout/fintech-go/pkg/money"
	"github.com/lightsgoout/fintech-go/pkg/testing/isolation"
	"testing"
)

func TestPaymentsService_CreateAccount(t *testing.T) {
	env := isolation.PrepareTest(t)
	defer env.Rollback()

	svc := NewPaymentsService(env.Tx)

	t.Run("create account OK", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		err := svc.CreateAccount(env.Ctx, "bob", money.NewNumericFromInt64(100), "USD")
		if err != nil {
			t.Error(err)
		}
	}))

	t.Run("disallow balance < 0", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		err := svc.CreateAccount(env.Ctx, "bob", money.NewNumericFromInt64(-1), "USD")
		if !errors.Is(err, service.ErrInsufficientFunds) {
			t.Errorf("expected ErrInsufficientFunds, got err=%v", err)
		}
	}))

	t.Run("duplicate account disallowed", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		err := svc.CreateAccount(env.Ctx, "bob", money.NewNumericFromInt64(60), "USD")
		if err != nil {
			t.Error(err)
		}
		err = svc.CreateAccount(env.Ctx, "bob", money.NewNumericFromInt64(60), "USD")
		if !errors.Is(err, service.ErrAccountAlreadyExists) {
			t.Errorf("expected ErrAccountAlreadyExists, got err=%v", err)
		}
	}))

	t.Run("incorrect currency", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		err := svc.CreateAccount(env.Ctx, "bob", money.NewNumericFromInt64(60), "UAH")
		if !errors.Is(err, service.ErrIncompatibleCurrency) {
			t.Errorf("expected ErrIncompatibleCurrency, got err=%v", err)
		}
	}))

	t.Run("unexpected db error", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		if _, err := env.Tx.Exec(`alter table account drop column currency`); err != nil {
			t.Error(err)
		}
		err := svc.CreateAccount(env.Ctx, "bob", money.NewNumericFromInt64(60), "USD")
		if _, ok := err.(service.ErrInternal); !ok {
			t.Fail()
		}
	}))

	t.Run("bad account id", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		err := svc.CreateAccount(env.Ctx, "", money.NewNumericFromInt64(60), "EUR")
		if !errors.Is(err, service.ErrBadAccountID) {
			t.Errorf("expected ErrBadAccountID, got err=%v", err)
		}
	}))
}
