package main

import (
	"log"
	"os"

	"github.com/lekss361/currencyservice/currency/internal/clients/currencyclient"
	"github.com/lekss361/currencyservice/currency/internal/db"
	"github.com/lekss361/currencyservice/currency/internal/repository"
	"github.com/robfig/cron/v3"
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

	ratesRepo := repository.NewRatesRepo(sqlDB)

	c := cron.New()

	_, err = c.AddFunc("@daily", func() {
		log.Println("Starting daily currency fetch…")

		ratesResp, err := currencyclient.GetCurs()
		if err != nil {
			log.Printf("error fetching rates: %v", err)
			return
		}

		if err := ratesRepo.Save(ratesResp.Date, ratesResp.Rub); err != nil {
			log.Printf("error saving rates: %v", err)
			return
		}

		log.Printf("Successfully saved rates for %s", ratesResp.Date.Format("2006-01-02"))
	})
	if err != nil {
		log.Fatalf("failed to schedule daily job: %v", err)
	}

	c.Start()
	log.Println("Cron scheduler started — will run once a day at midnight.")

	select {}
}
