package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"bdd-website/internal/database"
	"bdd-website/internal/middleware"
	"bdd-website/internal/models"
)

// GetActivities récupère la liste des activités
func GetActivities(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer les paramètres de pagination
		page, pageSize := getPagination(r)

		// Déterminer si on affiche uniquement les activités à venir
		upcomingOnly := true
		if showAll := r.URL.Query().Get("all"); showAll == "true" {
			upcomingOnly = false
		}

		// Récupérer l'ID utilisateur (optionnel)
		userID := middleware.GetUserID(r)

		// Récupérer les activités
		activities, total, err := database.GetActivities(db, page, pageSize, upcomingOnly, userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération des activités")
			return
		}

		// Répondre avec la liste des activités
		respondWithJSON(w, http.StatusOK, models.ActivitiesResponse{
			Activities: activities,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
		})
	}
}

// GetActivity récupère les détails d'une activité
func GetActivity(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID de l'activité
		activityID, err := getIDParam(r, "id")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "ID d'activité invalide")
			return
		}

		// Récupérer l'ID utilisateur (optionnel)
		userID := middleware.GetUserID(r)

		// Récupérer l'activité
		activity, err := database.GetActivity(db, activityID, userID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Activité non trouvée")
			return
		}

		// Répondre avec les détails de l'activité
		respondWithJSON(w, http.StatusOK, activity)
	}
}

// RegisterToActivity inscrit un utilisateur à une activité
func RegisterToActivity(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID utilisateur du contexte
		userID, ok := getRequiredUserID(w, r)
		if !ok {
			return
		}

		// Récupérer l'ID de l'activité
		activityID, err := getIDParam(r, "id")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "ID d'activité invalide")
			return
		}

		// Inscrire l'utilisateur à l'activité
		err = database.RegisterToActivity(db, userID, activityID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Répondre avec succès
		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Inscription réussie à l'activité",
		})
	}
}

// UnregisterFromActivity désinscrire un utilisateur d'une activité
func UnregisterFromActivity(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID utilisateur du contexte
		userID, ok := getRequiredUserID(w, r)
		if !ok {
			return
		}

		// Récupérer l'ID de l'activité
		activityID, err := getIDParam(r, "id")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "ID d'activité invalide")
			return
		}

		// Désinscrire l'utilisateur de l'activité
		err = database.UnregisterFromActivity(db, userID, activityID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Répondre avec succès
		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Désinscription réussie de l'activité",
		})
	}
}

// GetUserRegistrations récupère les activités auxquelles un utilisateur est inscrit
func GetUserRegistrations(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID utilisateur du contexte
		userID, ok := getRequiredUserID(w, r)
		if !ok {
			return
		}

		// Déterminer si on inclut l'historique
		includeHistory := false
		if history := r.URL.Query().Get("history"); history == "true" {
			includeHistory = true
		}

		// Récupérer les activités
		activities, err := database.GetUserRegistrations(db, userID, includeHistory)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération des inscriptions")
			return
		}

		// Répondre avec la liste des activités
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"activities": activities,
			"count":      len(activities),
		})
	}
}

// AdminCreateActivity permet à un administrateur de créer une activité
func AdminCreateActivity(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Décoder le corps de la requête
		var activityCreate models.ActivityCreate
		if err := json.NewDecoder(r.Body).Decode(&activityCreate); err != nil {
			respondWithError(w, http.StatusBadRequest, "Format de requête invalide")
			return
		}

		// Valider les données
		if activityCreate.Title == "" || activityCreate.Description == "" {
			respondWithError(w, http.StatusBadRequest, "Titre et description obligatoires")
			return
		}

		// Créer l'activité
		activityID, err := database.CreateActivity(db, activityCreate)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la création de l'activité")
			return
		}

		// Récupérer l'activité créée
		activity, err := database.GetActivity(db, activityID, 0)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération de l'activité")
			return
		}

		// Répondre avec l'activité créée
		respondWithJSON(w, http.StatusCreated, activity)
	}
}

// AdminUpdateActivity permet à un administrateur de mettre à jour une activité
func AdminUpdateActivity(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID de l'activité
		activityID, err := getIDParam(r, "id")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "ID d'activité invalide")
			return
		}

		// Décoder le corps de la requête
		var activityUpdate models.ActivityUpdate
		if err := json.NewDecoder(r.Body).Decode(&activityUpdate); err != nil {
			respondWithError(w, http.StatusBadRequest, "Format de requête invalide")
			return
		}

		// Valider les données
		if activityUpdate.Title == "" || activityUpdate.Description == "" {
			respondWithError(w, http.StatusBadRequest, "Titre et description obligatoires")
			return
		}

		// Mettre à jour l'activité
		err = database.UpdateActivity(db, activityID, activityUpdate)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Récupérer l'activité mise à jour
		activity, err := database.GetActivity(db, activityID, 0)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération de l'activité")
			return
		}

		// Répondre avec l'activité mise à jour
		respondWithJSON(w, http.StatusOK, activity)
	}
}

// AdminDeleteActivity permet à un administrateur de supprimer une activité
func AdminDeleteActivity(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID de l'activité
		activityID, err := getIDParam(r, "id")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "ID d'activité invalide")
			return
		}

		// Supprimer l'activité
		err = database.DeleteActivity(db, activityID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Répondre avec succès
		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Activité supprimée avec succès",
		})
	}
}
