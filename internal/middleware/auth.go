package middleware

import (
	"context"
	"net/http"
	"strings"

	"bdd-website/internal/utils"
)

// Clés de contexte pour les informations utilisateur
type contextKey string

const (
	UserIDKey  contextKey = "user_id"
	IsAdminKey contextKey = "is_admin"
)

// Auth est un middleware pour vérifier l'authentification par JWT
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extraire le token du header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			// Format du token: "Bearer <token>"
			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			tokenString := bearerToken[1]

			// Valider le token
			claims, err := utils.ValidateToken(tokenString, jwtSecret)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Ajouter les informations utilisateur au contexte de la requête
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, IsAdminKey, claims.IsAdmin)

			// Passer au gestionnaire suivant avec le contexte mis à jour
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID récupère l'ID utilisateur depuis le contexte
func GetUserID(r *http.Request) int64 {
	userID, ok := r.Context().Value(UserIDKey).(int64)
	if !ok {
		return 0
	}
	return userID
}

// IsAdmin vérifie si l'utilisateur est administrateur
func IsAdmin(r *http.Request) bool {
	isAdmin, ok := r.Context().Value(IsAdminKey).(bool)
	if !ok {
		return false
	}
	return isAdmin
}

// AdminOnly est un middleware qui vérifie si l'utilisateur est administrateur
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAdmin(r) {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
