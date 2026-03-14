package handlers

import (
	"database/sql"
	"net/http"

	"obelisk/internal/webui"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	// @AlexMcHugh1: note that we add the API routes first,
	// so they take precedence over the web UI routes.
	registerAPIRoutes(mux, db)
	registerWebRoutes(mux)

	return mux
}

func registerWebRoutes(mux *http.ServeMux) {
	web := webui.Handler()
	mux.Handle("/", web)
}
