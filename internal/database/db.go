package db

import (
	"database/sql"
	"fmt"
	_"github.com/lib/pq"
)

const (
	host	="localhost"
	port	= 5432
	user	= "devuser"
	password	= "devpassword"
	dbname	= "pdfvault"
	sslmode	= "disable"
)

// initialise the connection pool and return it for use in main.go
func InitDB() (*sql.DB, error) {
	// format the connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
		
	// open the connection pool
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// verify the connection is active
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to the database")
	return db, nil
}