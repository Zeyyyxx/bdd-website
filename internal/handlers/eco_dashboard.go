package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"bdd-website/internal/database"
	"bdd-website/internal/models"
)

// GetUserEcoPoints récupère les points écologiques d'un utilisateur
func GetUserEcoPoints(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID utilisateur du contexte
		userID, ok := getRequiredUserID(w, r)
		if !ok {
			return
		}

		// Récupérer les points
		points, totalPoints, err := database.GetUserEcoPoints(db, userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération des points")
			return
		}

		// Répondre avec les points
		respondWithJSON(w, http.StatusOK, models.EcoPointsResponse{
			Points:      points,
			TotalPoints: totalPoints,
		})
	}
}

// GetUserChallenges récupère les défis écologiques disponibles et en cours
func GetUserChallenges(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID utilisateur du contexte
		userID, ok := getRequiredUserID(w, r)
		if !ok {
			return
		}

		// Récupérer les défis actifs
		activeChallenges, err := database.GetChallenges(db, userID, true)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération des défis actifs")
			return
		}

		// Récupérer les défis auxquels l'utilisateur participe
		userChallenges, err := database.GetUserChallenges(db, userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération des défis utilisateur")
			return
		}

		// Filtrer les défis complétés
		var inProgressChallenges []models.Challenge
		var completedChallenges []models.Challenge

		for _, challenge := range userChallenges {
			if challenge.UserStatus == "completed" {
				completedChallenges = append(completedChallenges, challenge)
			} else {
				inProgressChallenges = append(inProgressChallenges, challenge)
			}
		}

		// Répondre avec les défis
		respondWithJSON(w, http.StatusOK, models.ChallengesResponse{
			ActiveChallenges:    activeChallenges,
			UserChallenges:      inProgressChallenges,
			CompletedChallenges: completedChallenges,
		})
	}
}

// JoinChallenge permet à un utilisateur de rejoindre un défi
func JoinChallenge(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID utilisateur du contexte
		userID, ok := getRequiredUserID(w, r)
		if !ok {
			return
		}

		// Récupérer l'ID du défi
		challengeID, err := getIDParam(r, "id")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "ID de défi invalide")
			return
		}

		// Rejoindre le défi
		err = database.JoinChallenge(db, userID, challengeID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Répondre avec succès
		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Vous avez rejoint le défi avec succès",
		})
	}
}

// CompleteChallenge permet à un utilisateur de marquer un défi comme terminé
func CompleteChallenge(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID utilisateur du contexte
		userID, ok := getRequiredUserID(w, r)
		if !ok {
			return
		}

		// Récupérer l'ID du défi
		challengeID, err := getIDParam(r, "id")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "ID de défi invalide")
			return
		}

		// Terminer le défi
		err = database.CompleteChallenge(db, userID, challengeID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Répondre avec succès
		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Défi terminé avec succès",
		})
	}
}

// GetUserBadges récupère les badges d'un utilisateur
func GetUserBadges(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID utilisateur du contexte
		userID, ok := getRequiredUserID(w, r)
		if !ok {
			return
		}

		// Récupérer les badges
		earnedBadges, availableBadges, err := database.GetUserBadges(db, userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération des badges")
			return
		}

		// Répondre avec les badges
		respondWithJSON(w, http.StatusOK, models.BadgesResponse{
			EarnedBadges:    earnedBadges,
			AvailableBadges: availableBadges,
		})
	}
}

// AdminCreateChallenge permet à un administrateur de créer un défi
func AdminCreateChallenge(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Décoder le corps de la requête
		var challengeCreate models.ChallengeCreate
		if err := json.NewDecoder(r.Body).Decode(&challengeCreate); err != nil {
			respondWithError(w, http.StatusBadRequest, "Format de requête invalide")
			return
		}

		// Valider les données
		if challengeCreate.Title == "" || challengeCreate.Description == "" || challengeCreate.Points <= 0 {
			respondWithError(w, http.StatusBadRequest, "Titre, description et points positifs obligatoires")
			return
		}

		// Créer le défi
		challengeID, err := database.CreateChallenge(db, challengeCreate)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la création du défi")
			return
		}

		// Répondre avec succès
		respondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"message": "Défi créé avec succès",
			"id":      challengeID,
		})
	}
}

// AdminUpdateChallenge permet à un administrateur de mettre à jour un défi
func AdminUpdateChallenge(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID du défi
		challengeID, err := getIDParam(r, "id")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "ID de défi invalide")
			return
		}

		// Décoder le corps de la requête
		var challengeUpdate models.ChallengeUpdate
		if err := json.NewDecoder(r.Body).Decode(&challengeUpdate); err != nil {
			respondWithError(w, http.StatusBadRequest, "Format de requête invalide")
			return
		}

		// Valider les données
		if challengeUpdate.Title == "" || challengeUpdate.Description == "" || challengeUpdate.Points <= 0 {
			respondWithError(w, http.StatusBadRequest, "Titre, description et points positifs obligatoires")
			return
		}

		// Mettre à jour le défi
		err = database.UpdateChallenge(db, challengeID, challengeUpdate)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Répondre avec succès
		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Défi mis à jour avec succès",
		})
	}
}

// AdminDeleteChallenge permet à un administrateur de supprimer un défi
func AdminDeleteChallenge(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID du défi
		challengeID, err := getIDParam(r, "id")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "ID de défi invalide")
			return
		}

		// Supprimer le défi
		err = database.DeleteChallenge(db, challengeID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Répondre avec succès
		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Défi supprimé avec succès",
		})
	}
}
