package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"bdd-website/internal/database"
)

// AdminGetStats récupère les statistiques pour le tableau de bord administrateur
func AdminGetStats(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer les statistiques
		stats, err := database.GetAdminStats(db)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la récupération des statistiques")
			return
		}

		// Répondre avec les statistiques
		respondWithJSON(w, http.StatusOK, stats)
	}
}

// AdminMarkContactMessageAsRead permet à un administrateur de marquer un message comme lu
func AdminMarkContactMessageAsRead(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID du message
		messageID, err := getIDParam(r, "id")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "ID de message invalide")
			return
		}

		// Marquer comme lu
		err = database.MarkContactMessageAsRead(db, messageID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Répondre avec succès
		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Message marqué comme lu",
		})
	}
}

// AdminGetContactMessage permet à un administrateur de récupérer un message spécifique
func AdminGetContactMessage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID du message
		messageID, err := getIDParam(r, "id")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "ID de message invalide")
			return
		}

		// Récupérer le message
		message, err := database.GetContactMessage(db, messageID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Message non trouvé")
			return
		}

		// Marquer comme lu si pas déjà lu
		if !message.IsRead {
			err = database.MarkContactMessageAsRead(db, messageID)
			if err != nil {
				// Ne pas échouer la requête si le marquage échoue
				// Juste enregistrer l'erreur (dans un vrai système, on utiliserait un logger)
			}
			message.IsRead = true
		}

		// Répondre avec le message
		respondWithJSON(w, http.StatusOK, message)
	}
}

// AdminUpdateUserAdmin met à jour le statut administrateur d'un utilisateur
func AdminUpdateUserAdmin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'ID de l'utilisateur
		userID, err := getIDParam(r, "id")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "ID utilisateur invalide")
			return
		}

		// Vérifier si la requête contient un champ is_admin
		var req struct {
			IsAdmin bool `json:"is_admin"`
		}

		if err := decodeJSONBody(r, &req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Format de requête invalide")
			return
		}

		// Mettre à jour le statut admin
		err = database.UpdateUserAdminStatus(db, userID, req.IsAdmin)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Erreur lors de la mise à jour du statut admin")
			return
		}

		// Répondre avec succès
		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Statut administrateur mis à jour avec succès",
		})
	}
}

// decodeJSONBody décode le corps JSON d'une requête dans une structure
func decodeJSONBody(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(dst)
}
