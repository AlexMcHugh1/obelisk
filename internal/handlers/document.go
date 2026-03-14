package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"obelisk/internal/database"
	"obelisk/internal/models"
	"obelisk/pkg/helpers"
	"os"
	"path/filepath"
	"strconv"
)

func UploadFile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		userIDStr := r.FormValue("userID")
		userID, _ := strconv.Atoi(userIDStr)

		file, handler, err := r.FormFile("myFile")
		if err != nil {
			http.Error(w, "File not found", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// @AlexMcHugh1: we use the UPLOAD_DIR var from the helpers package
		// to determine where to save the uploaded files.
		// This allows us to easily change the upload directory in one place
		// (like in the Dockerfile or docker-compose.yml)
		// without having to modify the code.
		dstPath := filepath.Join(helpers.UPLOAD_DIR, handler.Filename)
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
			UploaderID: userID,
		}

		if err := database.CreateDocument(db, doc); err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Uploaded and indexed: %s", handler.Filename)
	}
}

// Stream the requested PDF back to the client
func DownloadFile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, "Invalid document ID", http.StatusBadRequest)
			return
		}
		if id == 0 {
			http.Error(w, "Missing document ID", http.StatusBadRequest)
			return
		}

		// Fetches metadata from the DB to find the physical file path
		doc, err := database.GetDocumentByID(db, id)
		if err != nil {
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}

		// Sets headers to tell the browser it is a file download
		w.Header().Set("Content-Disposition", "attachment; filename="+doc.Title)
		http.ServeFile(w, r, doc.FilePath)
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

func ShareDocument(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docID := r.URL.Query().Get("doc_id")
		targetUserID := r.URL.Query().Get("target_user_id")

		if docID == "" || targetUserID == "" {
			http.Error(w, "Missing doc_id or target_user_id", http.StatusBadRequest)
			return
		}

		query := `INSERT INTO shared_access (document_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
		_, err := db.Exec(query, docID, targetUserID)
		if err != nil {
			http.Error(w, "Failed to share document", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Document %s shared with user %s", docID, targetUserID)
	}
}

// My documents tab
func MyDocumentsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		docs, err := database.GetUserDocuments(db, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(docs)
	}
}

// Shared tab
func SharedWithMeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		docs, err := database.GetSharedWithMeDocuments(db, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(docs)
	}
}
