package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"obelisk/internal/auth"
	"obelisk/internal/database"
	"obelisk/internal/models"
)

func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		hashed, _ := auth.HashPassword(creds.Password)
		newUser := models.User{
			Username:     creds.Username,
			PasswordHash: hashed,
			Role:         "user",
		}

		if err := database.CreateUser(db, newUser); err != nil {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "Account created successfully")
	}
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		user, err := database.GetUserByUsername(db, creds.Username)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		if !auth.CheckPasswordHash(creds.Password, user.PasswordHash) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// The frontend requires these specific keys
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id":  user.ID,
			"username": user.Username,
		})
	}
}
