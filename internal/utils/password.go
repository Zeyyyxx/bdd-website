package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Coût du hachage bcrypt (plus élevé = plus sécurisé mais plus lent)
const bcryptCost = 10

// HashPassword génère un hash bcrypt à partir d'un mot de passe en texte brut
func HashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("le mot de passe doit contenir au moins 6 caractères")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// CheckPasswordHash vérifie si un mot de passe correspond à un hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
