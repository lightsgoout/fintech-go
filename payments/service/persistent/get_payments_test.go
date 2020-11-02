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

func TestPaymentsService_GetPayments(t *testing.T) {
	env := isolation.PrepareTest(t)
	defer env.Rollback()

	svc := NewPaymentsService(env.Tx)

	t.Run("check account exists", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		res, err := svc.GetPayments(env.Ctx, "zzz")
		if !errors.Is(err, service.ErrAccountDoesNotExist) {
			t.Errorf("expected ErrAccountDoesNotExist, got result=%v, err=%v", res, err)
		}
	}))

	t.Run("get empty payments", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		err := svc.CreateAccount(env.Ctx, "bob", money.NewNumericFromInt64(100), "USD")
		if err != nil {
			t.Error(err)
		}
		res, err := svc.GetPayments(env.Ctx, "bob")
		if err != nil {
			t.Error(err)
		}
		if len(res) != 0 {
			t.Errorf("expected empty result, got len %d", len(res))
		}
	}))

	t.Run("get payments OK", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		// Create accounts
		const bob = entity.AccountID("bob")
		const alice = entity.AccountID("alice")
		for _, id := range [...]entity.AccountID{bob, alice} {
			err := svc.CreateAccount(env.Ctx, id, money.NewNumericFromInt64(100), "USD")
			if err != nil {
				t.Error(err)
			}
		}

		// Transfer some money from bob to alice
		_, err := svc.Transfer(env.Ctx, bob, alice, money.NewNumericFromInt64(50), "USD")
		if err != nil {
			t.Error(err)
		}
		// And from alice to bob
		_, err = svc.Transfer(env.Ctx, alice, bob, money.NewNumericFromInt64(35), "USD")
		if err != nil {
			t.Error(err)
		}

		// Check bob's payments
		bobPayments, err := svc.GetPayments(env.Ctx, bob)
		if err != nil {
			t.Error(err)
		}
		if len(bobPayments) != 2 {
			t.Error("incorrect bob payments count")
		}
		assert.Equal(t, bobPayments[0].Value.From, alice)
		assert.Equal(t, bobPayments[0].Value.To, bob)
		assert.Equal(t, bobPayments[0].Value.Amount, money.NewNumericFromInt64(35))
		assert.Equal(t, bobPayments[0].Value.Currency, money.Currency("USD"))
		assert.Equal(t, bobPayments[0].Value.Outgoing, false)
		assert.Equal(t, bobPayments[1].Value.From, bob)
		assert.Equal(t, bobPayments[1].Value.To, alice)
		assert.Equal(t, bobPayments[1].Value.Amount, money.NewNumericFromInt64(50))
		assert.Equal(t, bobPayments[1].Value.Currency, money.Currency("USD"))
		assert.Equal(t, bobPayments[1].Value.Outgoing, true)

		// Check alice's payments
		alicePayments, err := svc.GetPayments(env.Ctx, alice)
		if err != nil {
			t.Error(err)
		}
		if len(alicePayments) != 2 {
			t.Error("incorrect alice payments count")
		}
		assert.Equal(t, alicePayments[0].Value.From, alice)
		assert.Equal(t, alicePayments[0].Value.To, bob)
		assert.Equal(t, alicePayments[0].Value.Amount, money.NewNumericFromInt64(35))
		assert.Equal(t, alicePayments[0].Value.Currency, money.Currency("USD"))
		assert.Equal(t, alicePayments[0].Value.Outgoing, true)
		assert.Equal(t, alicePayments[1].Value.From, bob)
		assert.Equal(t, alicePayments[1].Value.To, alice)
		assert.Equal(t, alicePayments[1].Value.Amount, money.NewNumericFromInt64(50))
		assert.Equal(t, alicePayments[1].Value.Currency, money.Currency("USD"))
		assert.Equal(t, alicePayments[1].Value.Outgoing, false)
	}))
}
