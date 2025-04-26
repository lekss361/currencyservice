package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lekss361/currencyservice/currency/internal/clients/currencyclient"
	"github.com/lekss361/currencyservice/currency/internal/db"
	"github.com/lekss361/currencyservice/currency/internal/handler"
	"github.com/lekss361/currencyservice/currency/internal/repository"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime)

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	sqlDB, err := db.New(dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer sqlDB.Close()

	if err := db.Migrate(sqlDB); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	ratesRepo := repository.NewRatesRepo(sqlDB)

	ratesResp, err := currencyclient.GetCurs()
	if err != nil {
		log.Fatalf("error fetching currency rates: %v", err)
	}
	if err := ratesRepo.Save(ratesResp.Date, ratesResp.Rub); err != nil {
		log.Fatalf("error saving rates: %v", err)
	}
	fmt.Println("Курсы сохранены успешно")

	mux := http.NewServeMux()
	h := handler.New(ratesRepo)
	h.RegisterRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "12358"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting HTTP server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
