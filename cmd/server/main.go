package main

import (
	"fmt"
	"log"
	"obelisk/internal/database"
)

func main() {
	// Change 'database.InitDB()' to 'db.InitDB()'
	conn, err := db.InitDB() 
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Obelisk is online and connected to the database.")
}