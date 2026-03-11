package database

import (
	"database/sql"
	"fmt"
	"obelisk/internal/models"

	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "devuser"
	password = "devpassword"
	dbname   = "obelisk"
	sslmode  = "disable"
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

// Save file metadata to the database
func CreateDocument(db *sql.DB, doc models.Document) error {
	query := `INSERT INTO documents (title, file_path, uploader_id) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, doc.Title, doc.FilePath, doc.UploaderID)
	return err
}
func GetDocuments(db *sql.DB) ([]models.Document, error) {
	rows, err := db.Query("SELECT id, title, file_path, uploader_id, upload_date FROM documents")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	docs := []models.Document{}
	for rows.Next() {
		var d models.Document
		if err := rows.Scan(&d.ID, &d.Title, &d.FilePath, &d.UploaderID, &d.UploadDate); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}

// Find a single document's metadata by its database ID
func GetDocumentByID(db *sql.DB, id string) (models.Document, error) {
	var d models.Document
	query := `SELECT id, title, file_path, uploader_id, upload_date FROM documents WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&d.ID, &d.Title, &d.FilePath, &d.UploaderID, &d.UploadDate)
	return d, err
}

// Fetch files the user actually uploaded (My Documents)
func GetUserDocuments(db *sql.DB, userID string) ([]models.Document, error) {
	query := `SELECT id, title, file_path, uploader_id, upload_date FROM documents WHERE uploader_id = $1`
	return queryDocuments(db, query, userID)
}

// Fetch files others shared with this user (Shared Tab)
func GetSharedWithMeDocuments(db *sql.DB, userID string) ([]models.Document, error) {
	query := `
		SELECT d.id, d.title, d.file_path, d.uploader_id, d.upload_date 
		FROM documents d
		JOIN shared_access s ON d.id = s.document_id
		WHERE s.user_id = $1`
	return queryDocuments(db, query, userID)
}

func queryDocuments(db *sql.DB, query string, args ...interface{}) ([]models.Document, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	docs := []models.Document{}
	for rows.Next() {
		var d models.Document
		if err := rows.Scan(&d.ID, &d.Title, &d.FilePath, &d.UploaderID, &d.UploadDate); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}

// Find a user in the database by their unique username
func GetUserByUsername(db *sql.DB, username string) (models.User, error) {
	var u models.User
	query := `SELECT id, username, password_hash, role FROM users WHERE username = $1`
	err := db.QueryRow(query, username).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role)
	return u, err
}
