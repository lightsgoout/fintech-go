package get_payments

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/lightsgoout/fintech-go/payments/api/common"
	"github.com/lightsgoout/fintech-go/payments/entity"
	"github.com/lightsgoout/fintech-go/payments/service"
	"net/http"
	"time"
)

type getPaymentsRequest struct {
	AccountId entity.AccountID `json:"account_id"`
}

type outPayment struct {
	Id       entity.PaymentID `json:"id"`
	Time     time.Time        `json:"time"`
	From     entity.AccountID `json:"from"`
	To       entity.AccountID `json:"to"`
	Amount   string           `json:"amount"`
	Currency string           `json:"currency"`
	Outgoing bool             `json:"outgoing"`
}

type getPaymentsResponse struct {
	Payments []outPayment `json:"payments,omitempty"`
	Err      string       `json:"err,omitempty"`
}

func getPaymentsEndpoint(svc service.PaymentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getPaymentsRequest)
		payments, err := svc.GetPayments(
			ctx,
			req.AccountId,
		)
		if err != nil {
			return getPaymentsResponse{nil, err.Error()}, nil
		}

		outPayments := make([]outPayment, 0, len(payments))
		for _, p := range payments {
			outPayments = append(outPayments, outPayment{
				Id:       p.Id,
				Time:     p.Value.Time,
				From:     p.Value.From,
				To:       p.Value.To,
				Amount:   p.Value.Amount.String(),
				Currency: string(p.Value.Currency),
				Outgoing: p.Value.Outgoing,
			})
		}
		return getPaymentsResponse{outPayments, ""}, nil
	}
}

func decodeGetPaymentsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getPaymentsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func Server(svc service.PaymentsService) *httptransport.Server {
	return httptransport.NewServer(
		getPaymentsEndpoint(svc),
		decodeGetPaymentsRequest,
		common.EncodeResponse,
	)
}
