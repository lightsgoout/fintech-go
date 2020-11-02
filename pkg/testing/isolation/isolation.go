package isolation

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/lightsgoout/fintech-go/pkg/postgres"
	"testing"
)

// WrapInTransaction wraps test in transaction to ensure test atomicity.
func WrapInTransaction(tx *pg.Tx, f func(*testing.T)) func(t *testing.T) {
	return func(t *testing.T) {
		// go-pg doesn't support nested transaction, so do it old-fashioned way
		_, err := tx.Exec(`SAVEPOINT tx_001`)
		if err != nil {
			panic(err)
		}

		f(t)

		_, err = tx.Exec(`ROLLBACK TO tx_001`)
		if err != nil {
			panic(err)
		}
	}
}

// TestEnv contains isolated environment for a test suite
type TestEnv struct {
	// Tx is a transaction wrapping entire test suite
	Tx *pg.Tx

	// Rollback is a callback to rollback transaction
	// Should be deferred in each test
	Rollback func()

	Ctx context.Context
}

// PrepareTest setups postgres environment for a test suite
func PrepareTest(t *testing.T) TestEnv {
	db := postgres.NewPostgresFromEnv()
	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
	}

	return TestEnv{
		Ctx: context.Background(),
		Tx:  tx,
		Rollback: func() {
			err = tx.Rollback()
			if err != nil {
				t.Error(err)
			}
		},
	}
}
