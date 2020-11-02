package api

import (
	"github.com/lightsgoout/fintech-go/payments/api/create_account"
	"github.com/lightsgoout/fintech-go/payments/api/transfer"
	"github.com/lightsgoout/fintech-go/payments/service"
	"net/http"
)

func NewServer(svc service.PaymentsService, listenOn string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/account/create", create_account.Server(svc))
	mux.Handle("/transfer", transfer.Server(svc))
	return &http.Server{
		Addr:    listenOn,
		Handler: mux,
	}
}
