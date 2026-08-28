package main

import (
	"context"
	"github.com/jb843051627/indigo-vat/internal/httpapi"
	"github.com/jb843051627/indigo-vat/internal/service"
	"github.com/jb843051627/indigo-vat/internal/store"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	path := os.Getenv("INDIGO_VAT_DB")
	if path == "" {
		path = "./data/indigo-vat.db"
	}
	db, err := store.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	svc := service.New(db)
	svc.Start(ctx)
	defer svc.Close()
	server := &http.Server{Addr: ":8080", Handler: httpapi.New(svc).Handler()}
	go func() { <-ctx.Done(); _ = server.Shutdown(context.Background()) }()
	log.Println("indigo-vat listening on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
