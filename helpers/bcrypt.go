package helpers

import "golang.org/x/crypto/bcrypt"

func HashingPassword(p string) string {
	cost := 8
	password := []byte(p)
	hash, _ := bcrypt.GenerateFromPassword(password, cost)

	return string(hash)
}
