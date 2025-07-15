package utils

import (
	"crypto/sha1"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Funciones de autenticación

func ComparePasswords(storedHash string, loginPass string) error {
	byteHash := []byte(storedHash)
	loginHash := []byte(loginPass)
	//Compara las contraseñas
	err := bcrypt.CompareHashAndPassword(byteHash, loginHash)
	// Si las contraseñas no son iguales
	if err != nil {
		return err
	}

	return nil
}

func GeneratePassword(password string) string {
	binpwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(binpwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func HashPassword(password string) string {
	//Encriptar base de datos utilizando sha1
	hasher := sha1.New()
	hasher.Write([]byte(password))
	sha := fmt.Sprintf("%X", hasher.Sum(nil))
	return sha
}
