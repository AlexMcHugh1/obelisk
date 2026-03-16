package models

import "time"

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
}

type Document struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	FilePath   string    `json:"file_path"`
	UploaderID int       `json:"uploader_id"`
	UploadDate time.Time `json:"upload_date"`
}
