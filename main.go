package main

import (
	"flag"
	"fmt"
	"github.com/lightsgoout/fintech-go/payments/api"
	"github.com/lightsgoout/fintech-go/payments/service/persistent"
	"github.com/lightsgoout/fintech-go/pkg/postgres"
	"log"
	"net/http"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
	)
	flag.Parse()

	pg := postgres.NewPostgresFromEnv()
	svc := persistent.NewPaymentsService(pg)

	srv := http.Server{
		Addr:    *listen,
		Handler: api.NewAPIServer(svc),
	}
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Print(fmt.Errorf("failed to listen and serve: %w", err))
	}
}
