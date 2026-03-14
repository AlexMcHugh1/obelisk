package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"obelisk/internal/models"

	"github.com/lib/pq"
)

// @AlexMcHugh1: this error is used to indicate that a user
// with the same username already exists in the database.
var ErrUserExists = errors.New("user already exists")

const defaultDBPingTimeout = 5 * time.Second

func CreateUser(db *sql.DB, user models.User) error {
	const query = `
		INSERT INTO users (username, password_hash, role)
		VALUES ($1, $2, $3)
	`

	_, err := db.Exec(query, user.Username, user.PasswordHash, user.Role)
	if err == nil {
		return nil
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return ErrUserExists
	}

	return fmt.Errorf("create user %q: %w", user.Username, err)
}

func InitDB() (*sql.DB, error) {
	// @AlexMcHugh1: added the getenv helper function to read database connection
	// parameters from environment variables,
	host := getenv("DB_HOST", "127.0.0.1")
	port := getenv("DB_PORT", "5432")
	user := getenv("DB_USER", "devuser")
	password := getenv("DB_PASSWORD", "devpassword")
	name := getenv("DB_NAME", "obelisk")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	configurePool(db)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	const schema = `
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
		);
	`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("migrate schema: %w", err)
	}

	return nil
}

func CreateDocument(db *sql.DB, doc models.Document) error {
	const query = `
		INSERT INTO documents (title, file_path, uploader_id)
		VALUES ($1, $2, $3)
	`

	if _, err := db.Exec(query, doc.Title, doc.FilePath, doc.UploaderID); err != nil {
		return fmt.Errorf("create document %q: %w", doc.Title, err)
	}

	return nil
}

func GetDocuments(db *sql.DB) ([]models.Document, error) {
	const query = `
		SELECT id, title, file_path, uploader_id, upload_date
		FROM documents
		ORDER BY upload_date DESC
	`

	return queryDocuments(db, query)
}

func GetDocumentByID(db *sql.DB, id int) (models.Document, error) {
	const query = `
		SELECT id, title, file_path, uploader_id, upload_date
		FROM documents
		WHERE id = $1
	`

	var doc models.Document
	if err := db.QueryRow(query, id).Scan(
		&doc.ID,
		&doc.Title,
		&doc.FilePath,
		&doc.UploaderID,
		&doc.UploadDate,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Document{}, sql.ErrNoRows
		}
		return models.Document{}, fmt.Errorf("get document by id %d: %w", id, err)
	}

	return doc, nil
}

func GetUserDocuments(db *sql.DB, userID int) ([]models.Document, error) {
	const query = `
		SELECT id, title, file_path, uploader_id, upload_date
		FROM documents
		WHERE uploader_id = $1
		ORDER BY upload_date DESC
	`

	return queryDocuments(db, query, userID)
}

func GetSharedWithMeDocuments(db *sql.DB, userID int) ([]models.Document, error) {
	const query = `
		SELECT d.id, d.title, d.file_path, d.uploader_id, d.upload_date
		FROM documents d
		INNER JOIN shared_access s ON d.id = s.document_id
		WHERE s.user_id = $1
		ORDER BY d.upload_date DESC
	`

	return queryDocuments(db, query, userID)
}

func GetUserByUsername(db *sql.DB, username string) (models.User, error) {
	const query = `
		SELECT id, username, password_hash, role
		FROM users
		WHERE username = $1
	`

	var user models.User
	if err := db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, sql.ErrNoRows
		}
		return models.User{}, fmt.Errorf("get user by username %q: %w", username, err)
	}

	return user, nil
}

func queryDocuments(db *sql.DB, query string, args ...any) ([]models.Document, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query documents: %w", err)
	}
	defer rows.Close()

	var docs []models.Document
	for rows.Next() {
		var doc models.Document
		if err := rows.Scan(
			&doc.ID,
			&doc.Title,
			&doc.FilePath,
			&doc.UploaderID,
			&doc.UploadDate,
		); err != nil {
			return nil, fmt.Errorf("scan document row: %w", err)
		}
		docs = append(docs, doc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate document rows: %w", err)
	}

	return docs, nil
}

func configurePool(db *sql.DB) {
	db.SetMaxOpenConns(getenvInt("DB_MAX_OPEN_CONNS", 25))
	db.SetMaxIdleConns(getenvInt("DB_MAX_IDLE_CONNS", 25))
	db.SetConnMaxLifetime(getenvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute))
	db.SetConnMaxIdleTime(getenvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute))
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return n
}

func getenvDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return d
}
