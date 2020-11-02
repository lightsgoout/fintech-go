package api

import (
	"github.com/gorilla/mux"
	"github.com/lightsgoout/fintech-go/payments/api/create_account"
	"github.com/lightsgoout/fintech-go/payments/api/get_accounts"
	"github.com/lightsgoout/fintech-go/payments/api/get_payments"
	"github.com/lightsgoout/fintech-go/payments/api/transfer"
	"github.com/lightsgoout/fintech-go/payments/service"
	"net/http"
)

func NewAPIServer(svc service.PaymentsService) http.Handler {
	router := mux.NewRouter()
	router.Methods("POST").Path("/account/create").Handler(create_account.Server(svc))
	router.Methods("POST").Path("/transfer").Handler(transfer.Server(svc))
	router.Methods("POST").Path("/account/list").Handler(get_accounts.Server(svc))
	router.Methods("POST").Path("/payment/list").Handler(get_payments.Server(svc))
	return router
}
