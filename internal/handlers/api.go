package handlers

import (
	"database/sql"
	"net/http"
)

// @AlexMcHugh1: more idiomatic to have a single router.go
// file that sets up all the routes,
// rather than scattering them across multiple files.
// This way, you can easily see all the available endpoints
// in one place and manage them more effectively.
const apiV1 = "/api/v1"

func registerAPIRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc(apiV1+"/register", Register(db))
	mux.HandleFunc(apiV1+"/login", Login(db))
	mux.HandleFunc(apiV1+"/upload", UploadFile(db))
	mux.HandleFunc(apiV1+"/documents", ListFiles(db))
	mux.HandleFunc(apiV1+"/documents/my", MyDocumentsHandler(db))
	mux.HandleFunc(apiV1+"/documents/shared", SharedWithMeHandler(db))
	mux.HandleFunc(apiV1+"/documents/download", DownloadFile(db))
	mux.HandleFunc(apiV1+"/documents/share", ShareDocument(db))
}
