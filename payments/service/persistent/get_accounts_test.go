package persistent

import (
	"errors"
	"github.com/lightsgoout/fintech-go/payments/entity"
	"github.com/lightsgoout/fintech-go/payments/service"
	"github.com/lightsgoout/fintech-go/pkg/money"
	"github.com/lightsgoout/fintech-go/pkg/testing/isolation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPaymentsService_GetAccounts(t *testing.T) {
	env := isolation.PrepareTest(t)
	defer env.Rollback()

	svc := NewPaymentsService(env.Tx)

	t.Run("check currency", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		res, err := svc.GetAccounts(env.Ctx, money.NewCurrency("UAH"))
		if !errors.Is(err, service.ErrIncompatibleCurrency) {
			t.Errorf("expected ErrIncompatibleCurrency, got result=%v, err=%v", res, err)
		}
	}))

	t.Run("get accounts ok", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		// Create accounts
		const bob = entity.AccountID("bob")
		const alice = entity.AccountID("alice")
		err := svc.CreateAccount(env.Ctx, bob, money.NewNumericFromInt64(100), "USD")
		if err != nil {
			t.Error(err)
		}
		err = svc.CreateAccount(env.Ctx, alice, money.NewNumericFromInt64(100), "EUR")
		if err != nil {
			t.Error(err)
		}

		resultRUB, err := svc.GetAccounts(env.Ctx, money.NewCurrency("RUB"))
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, len(resultRUB), 0)

		resultUSD, err := svc.GetAccounts(env.Ctx, money.NewCurrency("USD"))
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, len(resultUSD), 1)
		assert.Equal(t, resultUSD[0], bob)

		resultEUR, err := svc.GetAccounts(env.Ctx, money.NewCurrency("EUR"))
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, len(resultEUR), 1)
		assert.Equal(t, resultEUR[0], alice)
	}))
}
