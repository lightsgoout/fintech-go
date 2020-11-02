package create_account

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/lightsgoout/fintech-go/payments/api/common"
	"github.com/lightsgoout/fintech-go/payments/entity"
	"github.com/lightsgoout/fintech-go/payments/service"
	"github.com/lightsgoout/fintech-go/pkg/money"
	"net/http"
)

type createAccountRequest struct {
	Id       entity.AccountID `json:"id"`
	Balance  money.Numeric    `json:"balance"`
	Currency money.Currency   `json:"currency"`
}

type createAccountResponse struct {
	Err string `json:"err,omitempty"`
}

func createAccountEndpoint(svc service.PaymentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createAccountRequest)
		err := svc.CreateAccount(ctx, req.Id, req.Balance, req.Currency)
		if err != nil {
			return createAccountResponse{err.Error()}, nil
		}
		return createAccountResponse{""}, err
	}
}

func decodeCreateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func Server(svc service.PaymentsService) *httptransport.Server {
	return httptransport.NewServer(
		createAccountEndpoint(svc),
		decodeCreateAccountRequest,
		common.EncodeResponse,
	)
}
