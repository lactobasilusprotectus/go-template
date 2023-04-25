package password

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword for encrypt password
func HashPassword(pass string) (hash string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		err = fmt.Errorf("Error password hash "+pass+", returns err: %+v", err)
		return
	}

	return string(hashedPassword), nil
}

// CheckPasswordHash for true or false
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
