package api

import (
	"encoding/json"
	"github.com/lightsgoout/fintech-go/payments/entity"
	"github.com/lightsgoout/fintech-go/payments/service/persistent"
	"github.com/lightsgoout/fintech-go/pkg/money"
	"github.com/lightsgoout/fintech-go/pkg/testing/isolation"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestServer_CreateAccount(t *testing.T) {
	env := isolation.PrepareTest(t)
	defer env.Rollback()

	svc := persistent.NewPaymentsService(env.Tx)
	srv := httptest.NewServer(NewAPIServer(svc))
	defer srv.Close()

	t.Run("create account ok", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		for _, testcase := range []struct {
			body, want string
		}{
			{
				body: `{"id":"bob","currency":"USD","balance":100.50}`,
				want: `{}`,
			},
		} {
			req, _ := http.NewRequest("POST", srv.URL+"/account/create", strings.NewReader(testcase.body))
			resp, _ := http.DefaultClient.Do(req)
			body, _ := ioutil.ReadAll(resp.Body)
			if want, have := testcase.want, strings.TrimSpace(string(body)); want != have {
				t.Errorf("%s : want %q, have %q", testcase.body, want, have)
			}
		}
	}))

	t.Run("lowercase currency", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		for _, testcase := range []struct {
			body, want string
		}{
			{
				body: `{"id":"bob","currency":"usd","balance":100.50}`,
				want: `{}`,
			},
		} {
			req, _ := http.NewRequest("POST", srv.URL+"/account/create", strings.NewReader(testcase.body))
			resp, _ := http.DefaultClient.Do(req)
			body, _ := ioutil.ReadAll(resp.Body)
			if want, have := testcase.want, strings.TrimSpace(string(body)); want != have {
				t.Errorf("%s : want %q, have %q", testcase.body, want, have)
			}
		}
	}))

	t.Run("unknown currency disallowed", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		for _, testcase := range []struct {
			body, want string
		}{
			{
				body: `{"id":"bob","currency":"UAH","balance":100.50}`,
				want: `{"err":"incompatible currency"}`,
			},
		} {
			req, _ := http.NewRequest("POST", srv.URL+"/account/create", strings.NewReader(testcase.body))
			resp, _ := http.DefaultClient.Do(req)
			body, _ := ioutil.ReadAll(resp.Body)
			if want, have := testcase.want, strings.TrimSpace(string(body)); want != have {
				t.Errorf("%s : want %q, have %q", testcase.body, want, have)
			}
		}
	}))

	t.Run("negative balance not allowed", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		for _, testcase := range []struct {
			body, want string
		}{
			{
				body: `{"id":"bob","currency":"UAH","balance":-600}`,
				want: `{"err":"insufficient funds"}`,
			},
		} {
			req, _ := http.NewRequest("POST", srv.URL+"/account/create", strings.NewReader(testcase.body))
			resp, _ := http.DefaultClient.Do(req)
			body, _ := ioutil.ReadAll(resp.Body)
			if want, have := testcase.want, strings.TrimSpace(string(body)); want != have {
				t.Errorf("%s : want %q, have %q", testcase.body, want, have)
			}
		}
	}))
}

func TestServer_Transfer(t *testing.T) {
	env := isolation.PrepareTest(t)
	defer env.Rollback()

	svc := persistent.NewPaymentsService(env.Tx)
	srv := httptest.NewServer(NewAPIServer(svc))
	defer srv.Close()

	// Create accounts
	const bob = entity.AccountID("bob")
	const alice = entity.AccountID("alice")
	for _, id := range [...]entity.AccountID{bob, alice} {
		err := svc.CreateAccount(env.Ctx, id, money.NewNumericFromInt64(100), "USD")
		if err != nil {
			t.Error(err)
		}
	}

	t.Run("transfer ok", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		for _, testcase := range []struct {
			body string
		}{
			{
				body: `{"from":"bob","to":"alice","currency":"USD","amount":30}`,
			},
			{
				body: `{"from":"alice","to":"bob","currency":"USD","amount":30}`,
			},
		} {
			req, _ := http.NewRequest("POST", srv.URL+"/transfer", strings.NewReader(testcase.body))
			resp, _ := http.DefaultClient.Do(req)
			body, _ := ioutil.ReadAll(resp.Body)

			var r struct {
				PaymentId int64 `json:"payment_id"`
			}
			err := json.Unmarshal(body, &r)
			if err != nil {
				t.Error(err)
			}
			assert.Greater(t, r.PaymentId, int64(0))
		}
	}))
}

func TestServer_GetAccounts(t *testing.T) {
	env := isolation.PrepareTest(t)
	defer env.Rollback()

	svc := persistent.NewPaymentsService(env.Tx)
	srv := httptest.NewServer(NewAPIServer(svc))
	defer srv.Close()

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

	t.Run("get accounts ok", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		for _, testcase := range []struct {
			body, want string
		}{
			{
				body: `{"currency":"USD"}`,
				want: `{"accounts":["bob"]}`,
			},
			{
				body: `{"currency":"EUR"}`,
				want: `{"accounts":["alice"]}`,
			},
			{
				body: `{"currency":"RUB"}`,
				want: `{}`,
			},
		} {
			req, _ := http.NewRequest("POST", srv.URL+"/account/list", strings.NewReader(testcase.body))
			resp, _ := http.DefaultClient.Do(req)
			body, _ := ioutil.ReadAll(resp.Body)
			if want, have := testcase.want, strings.TrimSpace(string(body)); want != have {
				t.Errorf("%s : want %q, have %q", testcase.body, want, have)
			}
		}
	}))

	t.Run("invalid currency", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		for _, testcase := range []struct {
			body, want string
		}{
			{
				body: `{"currency":"UAH"}`,
				want: `{"err":"incompatible currency"}`,
			},
		} {
			req, _ := http.NewRequest("POST", srv.URL+"/account/list", strings.NewReader(testcase.body))
			resp, _ := http.DefaultClient.Do(req)
			body, _ := ioutil.ReadAll(resp.Body)
			if want, have := testcase.want, strings.TrimSpace(string(body)); want != have {
				t.Errorf("%s : want %q, have %q", testcase.body, want, have)
			}
		}
	}))
}

func TestServer_GetPayments(t *testing.T) {
	env := isolation.PrepareTest(t)
	defer env.Rollback()

	svc := persistent.NewPaymentsService(env.Tx)
	srv := httptest.NewServer(NewAPIServer(svc))
	defer srv.Close()

	// Create accounts
	const bob = entity.AccountID("bob")
	const alice = entity.AccountID("alice")
	for _, id := range [...]entity.AccountID{bob, alice} {
		err := svc.CreateAccount(env.Ctx, id, money.NewNumericFromInt64(100), "USD")
		if err != nil {
			t.Error(err)
		}
	}

	paymentFromBob, err := svc.Transfer(env.Ctx, bob, alice, money.NewNumericFromInt64(30), money.NewCurrency("USD"))
	if err != nil {
		t.Error(err)
	}
	paymentFromAlice, err := svc.Transfer(env.Ctx, alice, bob, money.NewNumericFromInt64(25), money.NewCurrency("USD"))
	if err != nil {
		t.Error(err)
	}

	var r struct {
		Payments []struct {
			Id       entity.PaymentID `json:"id"`
			Time     time.Time        `json:"time"`
			From     entity.AccountID `json:"from"`
			To       entity.AccountID `json:"to"`
			Amount   string           `json:"amount"`
			Currency string           `json:"currency"`
			Outgoing bool             `json:"outgoing"`
		} `json:"payments"`
	}

	t.Run("get payments ok bob", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		req, _ := http.NewRequest("POST", srv.URL+"/payment/list", strings.NewReader(`{"account_id":"bob"}`))
		resp, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)

		err := json.Unmarshal(body, &r)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, r.Payments[0].Id, paymentFromAlice)
		assert.Equal(t, r.Payments[0].Outgoing, false)
		assert.Equal(t, r.Payments[0].Currency, "USD")
		assert.Equal(t, r.Payments[0].From, alice)
		assert.Equal(t, r.Payments[0].To, bob)
		assert.Equal(t, r.Payments[0].Amount, "25")
		assert.Equal(t, r.Payments[1].Id, paymentFromBob)
		assert.Equal(t, r.Payments[1].Outgoing, true)
		assert.Equal(t, r.Payments[1].Currency, "USD")
		assert.Equal(t, r.Payments[1].From, bob)
		assert.Equal(t, r.Payments[1].To, alice)
		assert.Equal(t, r.Payments[1].Amount, "30")
	}))

	t.Run("get payments ok alice", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		req, _ := http.NewRequest("POST", srv.URL+"/payment/list", strings.NewReader(`{"account_id":"alice"}`))
		resp, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)

		err := json.Unmarshal(body, &r)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, r.Payments[0].Id, paymentFromAlice)
		assert.Equal(t, r.Payments[0].Outgoing, true)
		assert.Equal(t, r.Payments[0].Currency, "USD")
		assert.Equal(t, r.Payments[0].From, alice)
		assert.Equal(t, r.Payments[0].To, bob)
		assert.Equal(t, r.Payments[0].Amount, "25")
		assert.Equal(t, r.Payments[1].Id, paymentFromBob)
		assert.Equal(t, r.Payments[1].Outgoing, false)
		assert.Equal(t, r.Payments[1].Currency, "USD")
		assert.Equal(t, r.Payments[1].From, bob)
		assert.Equal(t, r.Payments[1].To, alice)
		assert.Equal(t, r.Payments[1].Amount, "30")
	}))

	t.Run("incorrect account id", isolation.WrapInTransaction(env.Tx, func(t *testing.T) {
		req, _ := http.NewRequest("POST", srv.URL+"/payment/list", strings.NewReader(`{"account_id":""}`))
		resp, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, strings.TrimSpace(string(body)), `{"err":"bad account id"}`)
	}))
}
