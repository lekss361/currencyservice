package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lekss361/currencyservice/currency/internal/db"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime)

	dsnFlag := flag.String("dsn", "", "PostgreSQL DSN (e.g. \"postgres://user:pass@host:port/dbname?sslmode=disable\")")

	down := flag.Bool("down", false, "Run DOWN migrations (rollback). By default runs UP migrations.")
	flag.Parse()

	dsn := *dsnFlag
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	if dsn == "" {
		log.Fatal("DSN is not set: set -dsn flag or DATABASE_URL environment variable")
	}

	sqlDB, err := db.New(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer sqlDB.Close()

	if *down {
		log.Println("Rolling back migrations...")
		if err := db.MigrateDown(sqlDB); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		log.Println("Rollback completed")
	} else {
		log.Println("Applying migrations...")
		if err := db.Migrate(sqlDB); err != nil {
			log.Fatalf("Migrations failed: %v", err)
		}
		log.Println("Migrations applied")
	}

	fmt.Println("Done.")
}
