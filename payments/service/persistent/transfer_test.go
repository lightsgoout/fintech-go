package persistent

import (
	"errors"
	"github.com/lightsgoout/fintech-go/payments/entity"
	"github.com/lightsgoout/fintech-go/payments/service"
	"github.com/lightsgoout/fintech-go/pkg/money"
	"github.com/lightsgoout/fintech-go/pkg/testing/isolation"
	"testing"
)

func TestPaymentsService_Transfer(t *testing.T) {
	env := isolation.PrepareTest(t)
	defer env.Rollback()

	svc := NewPaymentsService(env.Tx)

	var (
		bob    entity.AccountID = "bob"
		bobEur entity.AccountID = "bob_eur"
		alice  entity.AccountID = "alice"
	)

	// create two accounts to test transfers between
	for _, id := range [...]entity.AccountID{bob, alice} {
		err := svc.CreateAccount(env.Ctx, id, money.NewNumericFromInt64(100), "USD")
		if err != nil {
			t.Error(err)
		}
	}
	err := svc.CreateAccount(env.Ctx, bobEur, money.NewNumericFromInt64(100), "EUR")
	if err != nil {
		t.Error(err)
	}

	t.Run("transfer 50 from bob to alice", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		_, err := svc.Transfer(env.Ctx, bob, alice, money.NewNumericFromInt64(50), "USD")
		if err != nil {
			t.Error(err)
		}
	}))

	t.Run("transfer 50 from alice to bob", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		_, err := svc.Transfer(env.Ctx, alice, bob, money.NewNumericFromInt64(50), "USD")
		if err != nil {
			t.Error(err)
		}
	}))

	t.Run("no transferring more than balance value", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		paymentId, err := svc.Transfer(env.Ctx, alice, bob, money.NewNumericFromInt64(9000), "USD")
		if !errors.Is(err, service.ErrInsufficientFunds) {
			t.Errorf("expected ErrInsufficientFunds, got paymentId=%v, err=%v", paymentId, err)
		}
	}))

	t.Run("only same currency allowed bob-side", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		paymentId, err := svc.Transfer(env.Ctx, alice, bob, money.NewNumericFromInt64(9000), "EUR")
		if !errors.Is(err, service.ErrIncompatibleCurrency) {
			t.Errorf("expected ErrIncompatibleCurrency, got paymentId=%v, err=%v", paymentId, err)
		}
	}))

	t.Run("only same currency allowed alice-side", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		paymentId, err := svc.Transfer(env.Ctx, bobEur, alice, money.NewNumericFromInt64(9000), "EUR")
		if !errors.Is(err, service.ErrIncompatibleCurrency) {
			t.Errorf("expected ErrIncompatibleCurrency, got paymentId=%v, err=%v", paymentId, err)
		}
	}))

	t.Run("check account exists", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		paymentId, err := svc.Transfer(env.Ctx, "abc", bob, money.NewNumericFromInt64(9000), "USD")
		if !errors.Is(err, service.ErrAccountDoesNotExist) {
			t.Errorf("expected ErrAccountDoesNotExist, got paymentId=%v, err=%v", paymentId, err)
		}
	}))

	t.Run("disallow transfer to the same account", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		paymentId, err := svc.Transfer(env.Ctx, bob, bob, money.NewNumericFromInt64(20), "USD")
		if !errors.Is(err, service.ErrBadTransferTarget) {
			t.Errorf("expected ErrBadTransferTarget, got paymentId=%v, err=%v", paymentId, err)
		}
	}))
}
