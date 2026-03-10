package database

import (
	"database/sql"
	"fmt"
	"obelisk/internal/models"
	_ "github.com/lib/pq"
)
const (
	host	="127.0.0.1"
	port	= 5432
	user	= "devuser"
	password	= "devpassword"
	dbname	= "obelisk"
	sslmode	= "disable"
)

// Insert a new user into the database
func CreateUser(db *sql.DB, user models.User) error {
	query := `INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, user.Username, user.PasswordHash, user.Role)
	return err
}

// Initialise the connection pool and return it for use in main.go
func InitDB() (*sql.DB, error) {
	// format the connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
		
	// Open the connection pool
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Verify the connection is active
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to the database")
	return db, nil
}

// Create the database tables if they do not exist
func Migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS documents (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		file_path TEXT NOT NULL,
		uploader_id INTEGER REFERENCES users(id),
		upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS shared_access (
		document_id INTEGER REFERENCES documents(id) ON DELETE CASCADE,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		shared_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (document_id, user_id)
	);`

	_, err := db.Exec(schema)
	return err
}


