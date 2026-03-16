package main

import (
	"fmt"
	"log"
	"net/http"
	"obelisk/internal/auth"
	"obelisk/internal/database"
	"obelisk/internal/handlers"
	"obelisk/internal/models"
	"os"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	hashedPwd, _ := auth.HashPassword("admin123")
	testUser := models.User{
		Username:     "admin",
		PasswordHash: hashedPwd,
		Role:         "admin",
	}

	err = database.CreateUser(db, testUser)
	if err != nil {
		fmt.Println("Note: Admin user was not created (likely already exists)")
	} else {
		fmt.Println("Admin user created successfully!")
	}

	rows, _ := db.Query("SELECT id, username FROM users")
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Printf("USER ID: %d | USERNAME: %s\n", id, name)
	}

	http.HandleFunc("/upload", handlers.UploadFile(db))
	http.HandleFunc("/list", handlers.ListFiles(db))
	http.HandleFunc("/my-docs", handlers.MyDocumentsHandler(db))
	http.HandleFunc("/shared-docs", handlers.SharedWithMeHandler(db))
	http.HandleFunc("/register", handlers.Register(db))
	http.HandleFunc("/login", handlers.Login(db))
	http.HandleFunc("/download", handlers.DownloadFile(db))
	http.HandleFunc("/share", handlers.ShareDocument(db))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	fmt.Printf("Obelisk is online. Server starting on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, enableCORS(http.DefaultServeMux)))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
