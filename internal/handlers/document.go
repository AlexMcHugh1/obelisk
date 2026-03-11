package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"obelisk/internal/database"
	"obelisk/internal/models"
	"os"
	"path/filepath"
)

func UploadFile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			http.Error(w, "File not found", http.StatusBadRequest)
			return
		}
		defer file.Close()

		dstPath := filepath.Join("uploads", handler.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Save failed", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Copy failed", http.StatusInternalServerError)
			return
		}

		doc := models.Document{
			Title:      handler.Filename,
			FilePath:   dstPath,
			UploaderID: 1,
		}

		if err := database.CreateDocument(db, doc); err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Uploaded and indexed: %s", handler.Filename)
	}
}

func ListFiles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docs, err := database.GetDocuments(db)
		if err != nil {
			http.Error(w, "Failed to fetch documents", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(docs)
	}
}
