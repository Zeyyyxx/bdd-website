package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"bdd-website/internal/database"
	"bdd-website/internal/models"
)

// SubmitContactForm traite l'envoi d'un formulaire de contact
func SubmitContactForm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Décoder le corps de la requête
		var contactCreate models.ContactMessageCreate
		if err := json.NewDecoder(r.Body).Decode(&contactCreate); err != nil {
			respondWithError(w, http.StatusBadRequest, "Format de requête invalide")
			return
		}

		// Valider les données
		if contactCreate.Name == "" || contactCreate.Email == "" ||
			contactCreate.Subject == "" || contactCreate.Message == "" {
			respondWithError(w, http.StatusBadRequest, "Tous les champs sont obligatoires")
			return
		}

		// Enregistrer le message
		messageID, err := database.CreateContactMessage(db, contactCreate)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de l'enregistrement du message")
			return
		}

		// Répondre avec succès
		respondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"message": "Message envoyé avec succès",
			"id":      messageID,
		})
	}
}

// AdminGetContactMessages permet à un administrateur de récupérer les messages de contact
func AdminGetContactMessages(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer les paramètres de pagination
		page, pageSize := getPagination(r)

		// Filtrer sur les messages non lus
		unreadOnly := false
		if unread := r.URL.Query().Get("unread"); unread == "true" {
			unreadOnly = true
		}

		// Récupérer les messages
		messages, total, unreadCount, err := database.GetContactMessages(db, page, pageSize, unreadOnly)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération des messages")
			return
		}

		// Répondre avec la liste des messages
		respondWithJSON(w, http.StatusOK, models.ContactMessagesResponse{
			Messages: messages,
			Total:    total,
			Unread:   unreadCount,
			Page:     page,
			PageSize: pageSize,
		})
	}
}

// AdminDeleteContactMessage permet à un administrateur de supprimer un message
func AdminDeleteContactMessage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID du message
		messageID, err := getIDParam(r, "id")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "ID de message invalide")
			return
		}

		// Supprimer le message
		err = database.DeleteContactMessage(db, messageID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Répondre avec succès
		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Message supprimé avec succès",
		})
	}
}
