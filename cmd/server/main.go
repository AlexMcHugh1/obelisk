package main

import (
	"fmt"
	"log"
	"obelisk/internal/auth"
	"obelisk/internal/database"
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
		log.Fatalf("Migration failed: %v", err)
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
		fmt.Println("Note: Admin user was not created (likely already exists)")
	} else {
		fmt.Println("Admin user created successfully!")
	}

	fmt.Println("Obelisk is online and connected to the database.")
}