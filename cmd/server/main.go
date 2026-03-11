package main

import (
	"fmt"
	"log"
	"net/http"
	"obelisk/internal/auth"
	"obelisk/internal/database"
	"obelisk/internal/handlers"
	"obelisk/internal/models"
)

func main() {
	// Initialize Connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Run Migrations (Create Tables)
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	// Create a Test Admin User
	hashedPwd, _ := auth.HashPassword("admin123")
	testUser := models.User{
		Username:     "admin",
		PasswordHash: hashedPwd,
		Role:         "admin",
	}

	err = database.CreateUser(db, testUser)
	if err != nil {
		// Unique constraint will prevent duplicates on subsequent runs
		fmt.Println("Note: Admin user was not created (likely already exists)")
	} else {
		fmt.Println("Admin user created successfully!")
	}

	// Set up Routes
	// Passing db into the handler allows it to record metadata
	http.HandleFunc("/upload", handlers.UploadFile(db))

	// Start Server
	fmt.Println("Obelisk is online. Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
