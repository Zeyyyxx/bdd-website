package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"bdd-website/internal/models"
)

// Claims représente les données encodées dans le JWT
type Claims struct {
	UserID  int64 `json:"user_id"`
	IsAdmin bool  `json:"is_admin"`
	jwt.RegisteredClaims
}

// GenerateToken crée un nouveau JWT pour l'utilisateur
func GenerateToken(user *models.User, secret string, expirationHours int) (string, error) {
	// Définir la durée d'expiration
	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	// Créer les claims
	claims := &Claims{
		UserID:  user.ID,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "bdd-website",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	// Créer le token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signer le token avec la clé secrète
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken vérifie si un token est valide et retourne les claims
func ValidateToken(tokenString string, secret string) (*Claims, error) {
	// Parse le token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Vérifier l'algorithme de signature
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Extraire et retourner les claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token invalide")
}
