package main

import (
	"context"
	"log"
	"os"

	"github.com/joaoprofile/gofi-sdk-go/base/environment"
	"github.com/joaoprofile/gofi-sdk-go/gofi"
	"github.com/joaoprofile/gofi-sdk-go/gofi/component/database"
	"github.com/joaoprofile/gofi-sdk-go/gofi/component/httpserver"
	"github.com/joaoprofile/gofi-sdk-go/netx/httpx"
	_ "github.com/joaoprofile/gofi-sdk-go/sqln/driver/postgres" // DATABASE_DRIVER=postgres

	"github.com/Arthur-155/art-books/domain/contas/handler"
	"github.com/Arthur-155/art-books/domain/contas/worker"
)

const (
	serviceName  = "art-books"
	httpPort     = ":8080"
	adminCommand = "admin"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	if len(os.Args) > 1 && os.Args[1] == adminCommand {
		return runAdmin(ctx, os.Args[2:])
	}
	cfg, err := LoadConfig(ctx, environment.Instance())
	if err != nil {
		return err
	}
	contas, tokens, err := buildContasService(cfg)
	if err != nil {
		return err
	}
	server := httpserver.New(httpPort, &httpx.WSConfig{
		AllowedOrigins: cfg.AllowedOrigins,
		Health:         &httpx.HealthConfig{},
	})
	app, err := gofi.New(serviceName).
		With(database.New(), worker.NewPurgeWorker(contas, worker.DefaultPurgeInterval), server).
		Build()
	if err != nil {
		return err
	}
	server.UseAuth(handler.Middleware(tokens)).
		Handlers(handler.NewContasHandler(contas), handler.NewAuthHandler(contas))
	return app.ListenAndServe()
}
