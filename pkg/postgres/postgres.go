package postgres

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"os"
)

// NewPostgresFromEnv returns new connection to Postgres using credentials from env
func NewPostgresFromEnv() *pg.DB {
	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgUser := os.Getenv("POSTGRES_USER")
	pgName := os.Getenv("POSTGRES_DB")
	pgPass := os.Getenv("POSTGRES_PASSWORD")

	return pg.Connect(&pg.Options{
		User:     pgUser,
		Database: pgName,
		Password: pgPass,
		Addr:     fmt.Sprintf("%s:%s", pgHost, pgPort),
	})
}

// Database is an interface conforming to both pg.DB and pg.Tx, necessary for isolated tests
type Database interface {
	ExecContext(c context.Context, query interface{}, params ...interface{}) (pg.Result, error)
	RunInTransaction(ctx context.Context, fn func(*pg.Tx) error) error
	QueryOneContext(c context.Context, model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
	QueryContext(c context.Context, model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
}

// NestedRunInTransaction runs func f in transaction, with support of nested transactions.
// go-pg doesn't support nested transactions when using RunInTransaction,
// so it poses a problem when running tests which are already inside a transaction.
// NOTE: supports only one level of nesting (enough for the purpose of this project)
func NestedRunInTransaction(ctx context.Context, db Database, f func(tx Database) error) error {
	// Check if we're already in a transaction
	// Kinda ugly :(
	switch db.(type) {
	case *pg.Tx:
		// We're inside a transaction, so need to use savepoints
		_, err := db.ExecContext(ctx, `SAVEPOINT tx_002`)
		if err != nil {
			return err
		}
		return f(db)
	case *pg.DB:
		// No transaction yet
		return (db.(*pg.DB)).RunInTransaction(ctx, func(tx *pg.Tx) error {
			return f(db)
		})
	default:
		// Must never get here
		panic("NestedRunInTransaction unexpected db type")
	}
	return nil
}
