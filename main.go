package main

import (
	"flag"
	"github.com/v1k45/pastepass/config"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/v1k45/pastepass/db"
	"github.com/v1k45/pastepass/web"
)

func main() {
	flag.StringVar(&config.ServerAddr, "server-addr", config.ServerAddr, "The server address to listen on")
	flag.StringVar(&config.AppName, "app-name", config.AppName, "The name of the application (e.g. ACME PastePass)")
	flag.StringVar(&config.DBPath, "db-path", config.DBPath, "The path to the database file")
	flag.BoolVar(&config.ResetDB, "reset-db", config.ResetDB, "Reset the database on startup")
	flag.Parse()

	// Open the database
	boltdb, err := db.NewDB(config.DBPath, config.ResetDB)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	go boltdb.DeleteExpiredPeriodically(time.Minute * 5)

	slog.Info("starting_server", "server_addr", config.ServerAddr, "app_name", config.AppName, "db_name", config.DBPath)

	// Start the web server
	handler := web.NewHandler(boltdb)
	http.ListenAndServe(config.ServerAddr, handler.Router())
}
