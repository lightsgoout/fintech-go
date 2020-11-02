package transfer

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

type transferRequest struct {
	From     entity.AccountID `json:"from"`
	To       entity.AccountID `json:"to"`
	Amount   money.Numeric    `json:"amount"`
	Currency money.Currency   `json:"currency"`
}

type transferResponse struct {
	PaymentId entity.PaymentID `json:"payment_id,omitempty"`
	Err       string           `json:"err,omitempty"`
}

func transferEndpoint(svc service.PaymentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transferRequest)
		paymentId, err := svc.Transfer(ctx, req.From, req.To, req.Amount, req.Currency)
		if err != nil {
			return transferResponse{0, err.Error()}, nil
		}
		return transferResponse{paymentId, ""}, nil
	}
}

func decodeTransferRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request transferRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func Server(svc service.PaymentsService) *httptransport.Server {
	return httptransport.NewServer(
		transferEndpoint(svc),
		decodeTransferRequest,
		common.EncodeResponse,
	)
}
