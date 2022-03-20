package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSaltPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func CheckPassword(pwd string, hash string) bool {
	if hash == "" {
		// if user doesn't hava pass / hash in DB
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}
