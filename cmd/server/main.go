package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"obelisk/internal/database"
	"obelisk/internal/handlers"
	"obelisk/pkg/helpers"
)

const addr = ":8080"

func main() {
	logger := log.Default()

	db, err := database.InitDB()
	if err != nil {
		logger.Fatalf("init database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Printf("close database: %v", err)
		}
	}()

	if err := database.Migrate(db); err != nil {
		logger.Fatalf("migrate database: %v", err)
	}

	if os.Getenv("UPLOAD_DIR") == "" {
		os.Setenv("UPLOAD_DIR", "./uploads")
	} else {
		// Ensure the directory exists
		if err := os.MkdirAll(os.Getenv("UPLOAD_DIR"), os.ModePerm); err != nil {
			panic("Failed to create upload directory: " + err.Error())
		}
	}
	helpers.UPLOAD_DIR = os.Getenv("UPLOAD_DIR")

	mux := handlers.NewRouter(db)

	srv := &http.Server{
		Addr:         addr,
		Handler:      withCORS(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Printf("Obelisk listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatalf("server stopped: %v", err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Set("Access-Control-Allow-Origin", "*")
		headers.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		headers.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
