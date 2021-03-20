package main

import (
	"context"
	"github.com/artrey/payments/cmd/service/app"
	"github.com/artrey/payments/pkg/business"
	"github.com/artrey/payments/pkg/security"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	defaultHost = "0.0.0.0"
	defaultPort = "9999"
	defaultDSN  = "postgres://app:pass@localhost:5432/db"
)

func main() {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}

	dsn, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		dsn = defaultDSN
	}

	if err := execute(net.JoinHostPort(host, port), dsn); err != nil {
		os.Exit(1)
	}
}

func execute(addr string, dsn string) error {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Println(err)
		return err
	}
	defer pool.Close()

	securitySvc := security.NewService(pool)
	businessSvc := business.NewService(pool)
	router := chi.NewRouter()
	application := app.NewServer(securitySvc, businessSvc, router)
	err = application.Init()
	if err != nil {
		log.Println(err)
		return err
	}

	server := &http.Server{
		Addr:    addr,
		Handler: application,
	}
	return server.ListenAndServe()
}
