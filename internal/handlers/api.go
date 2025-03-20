package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"bdd-website/internal/database"
	"bdd-website/internal/middleware"
)

// APIHealth vérifie l'état de l'API
func APIHealth(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// APISearchUsers permet de rechercher des utilisateurs (admin uniquement)
func APISearchUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Vérifier si l'utilisateur est admin
		if !middleware.IsAdmin(r) {
			respondWithError(w, http.StatusForbidden, "Accès interdit")
			return
		}

		// Récupérer le terme de recherche
		query := r.URL.Query().Get("q")
		if query == "" {
			respondWithError(w, http.StatusBadRequest, "Terme de recherche requis")
			return
		}

		// Récupérer la limite
		limit := 10
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
				// Limiter à un maximum raisonnable
				if limit > 50 {
					limit = 50
				}
			}
		}

		// Effectuer la recherche
		users, err := database.SearchUsers(db, query, limit)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la recherche des utilisateurs")
			return
		}

		// Répondre avec les résultats
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"users": users,
			"count": len(users),
		})
	}
}
