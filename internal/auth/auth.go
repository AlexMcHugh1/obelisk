package auth

import "golang.org/x/crypto/bcrypt"

// Turn plain text password into a secure hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}