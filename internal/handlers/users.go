package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"bdd-website/internal/database"
	"bdd-website/internal/models"
)

// GetUserProfile récupère le profil de l'utilisateur actuellement connecté
func GetUserProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID utilisateur du contexte
		userID, ok := getRequiredUserID(w, r)
		if !ok {
			return
		}

		// Récupérer le profil
		profile, err := database.GetUserProfile(db, userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération du profil")
			return
		}

		// Répondre avec le profil
		respondWithJSON(w, http.StatusOK, profile)
	}
}

// UpdateUserProfile met à jour le profil de l'utilisateur
func UpdateUserProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID utilisateur du contexte
		userID, ok := getRequiredUserID(w, r)
		if !ok {
			return
		}

		// Décoder le corps de la requête
		var profileUpdate models.UserProfileUpdate
		if err := json.NewDecoder(r.Body).Decode(&profileUpdate); err != nil {
			respondWithError(w, http.StatusBadRequest, "Format de requête invalide")
			return
		}

		// Valider les données (au moins un champ à mettre à jour)
		if profileUpdate.Username == "" && profileUpdate.Email == "" && profileUpdate.Password == "" {
			respondWithError(w, http.StatusBadRequest, "Aucun champ à mettre à jour")
			return
		}

		// Mettre à jour le profil
		err := database.UpdateUserProfile(db, userID, profileUpdate)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Récupérer le profil mis à jour
		profile, err := database.GetUserProfile(db, userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération du profil")
			return
		}

		// Répondre avec le profil mis à jour
		respondWithJSON(w, http.StatusOK, profile)
	}
}

// AdminGetUsers permet à un administrateur de récupérer la liste des utilisateurs
func AdminGetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer les paramètres de pagination
		page, pageSize := getPagination(r)

		// Récupérer les utilisateurs
		users, total, err := database.GetAllUsers(db, page, pageSize)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération des utilisateurs")
			return
		}

		// Répondre avec la liste des utilisateurs
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"users":     users,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		})
	}
}
