package main

import (
	"log"
	"net/http"
	"time"

	"github.com/v1k45/paste/db"
	"github.com/v1k45/paste/web"
)

func main() {
	// Open the database
	boltdb, err := db.NewDB("pastes.boltdb")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	go boltdb.DeleteExpiredPeriodically(time.Minute * 5)

	// Start the web server
	handler := web.NewHandler(boltdb)
	http.ListenAndServe(":8080", handler.Router())
}
