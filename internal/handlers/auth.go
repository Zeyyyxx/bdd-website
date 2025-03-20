package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"bdd-website/internal/database"
	"bdd-website/internal/models"
	"bdd-website/internal/utils"
)

// Register gère l'inscription d'un nouvel utilisateur
func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Décoder le corps de la requête
		var userRegister models.UserRegister
		if err := json.NewDecoder(r.Body).Decode(&userRegister); err != nil {
			respondWithError(w, http.StatusBadRequest, "Format de requête invalide")
			return
		}

		// Valider les données
		if userRegister.Email == "" || userRegister.Username == "" || userRegister.Password == "" {
			respondWithError(w, http.StatusBadRequest, "Tous les champs sont obligatoires")
			return
		}

		// Créer l'utilisateur
		userID, err := database.CreateUser(db, userRegister)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Récupérer l'utilisateur créé
		user, err := database.GetUserByID(db, userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération du profil")
			return
		}

		// Répondre avec succès
		respondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"message":  "Inscription réussie",
			"user_id":  userID,
			"email":    user.Email,
			"username": user.Username,
		})
	}
}

// Login gère la connexion d'un utilisateur
func Login(db *sql.DB, jwtSecret string, jwtExpirationHours int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Décoder le corps de la requête
		var userLogin models.UserLogin
		if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
			respondWithError(w, http.StatusBadRequest, "Format de requête invalide")
			return
		}

		// Valider les données
		if userLogin.Email == "" || userLogin.Password == "" {
			respondWithError(w, http.StatusBadRequest, "Email et mot de passe requis")
			return
		}

		// Récupérer l'utilisateur
		user, err := database.GetUserByEmail(db, userLogin.Email)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Email ou mot de passe incorrect")
			return
		}

		// Vérifier le mot de passe
		if !utils.CheckPasswordHash(userLogin.Password, user.Password) {
			respondWithError(w, http.StatusUnauthorized, "Email ou mot de passe incorrect")
			return
		}

		// Générer le token JWT
		token, err := utils.GenerateToken(user, jwtSecret, jwtExpirationHours)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la génération du token")
			return
		}

		// Récupérer le profil complet
		profile, err := database.GetUserProfile(db, user.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération du profil")
			return
		}

		// Répondre avec le token et les informations utilisateur
		respondWithJSON(w, http.StatusOK, models.UserResponse{
			User:  *profile,
			Token: token,
		})
	}
}
