package get_accounts

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

type getAccountsRequest struct {
	Currency string `json:"currency"`
}

type getAccountsResponse struct {
	Accounts []entity.AccountID `json:"accounts,omitempty"`
	Err      string             `json:"err,omitempty"`
}

func getAccountsEndpoint(svc service.PaymentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAccountsRequest)
		accs, err := svc.GetAccounts(
			ctx,
			money.NewCurrency(req.Currency),
		)
		if err != nil {
			return getAccountsResponse{nil, err.Error()}, nil
		}
		return getAccountsResponse{accs, ""}, nil
	}
}

func decodeGetAccountsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getAccountsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func Server(svc service.PaymentsService) *httptransport.Server {
	return httptransport.NewServer(
		getAccountsEndpoint(svc),
		decodeGetAccountsRequest,
		common.EncodeResponse,
	)
}
