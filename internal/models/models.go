package models

import "time"

// create structs
type User struct {
	ID	int	`json:"id"`
	Username string	`json:"username"`
	PasswordHash string	`json:"-"` // hide from JSON
	Role	string	`json:"role"`
}

type Document struct {
	ID	int	`json:"id"`
	Title	string	`json:"title"`
	FilePath	string	`json:"file_path"`
	UploaderID	int	`json:"uploader_id"`
	UploadDate	time.Time	`json:"upload_date"`
}