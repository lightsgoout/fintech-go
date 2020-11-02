package persistent

import (
	"fmt"
	"github.com/lightsgoout/fintech-go/payments/service"
	"github.com/lightsgoout/fintech-go/pkg/postgres"
)

// PaymentsService implements service.PaymentsService interface by using persistent storage.
//
// NOTE: pg can be pg.DB (in production) or pg.Tx (in tests), because we must isolate tests in transactions.
type PaymentsService struct {
	pg postgres.Database
}

// NewPaymentsService returns new PaymentsService with Postgres connection.
func NewPaymentsService(pg postgres.Database) PaymentsService {
	return PaymentsService{
		pg: pg,
	}
}

func NewInternalErrorFromDBError(err error) service.ErrInternal {
	return service.NewErrInternal(fmt.Errorf("database error: %w", err))
}
