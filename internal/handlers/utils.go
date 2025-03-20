package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"bdd-website/internal/middleware"
)

// Pagination par défaut
const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// respondWithError envoie une réponse d'erreur JSON
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON envoie une réponse JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Convertir le payload en JSON
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Erreur lors de la génération de la réponse JSON"}`))
		return
	}

	// Définir les headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// getPagination extrait les paramètres de pagination de la requête
func getPagination(r *http.Request) (page, pageSize int) {
	// Valeurs par défaut
	page = DefaultPage
	pageSize = DefaultPageSize

	// Récupérer les paramètres de requête
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
			// Limiter la taille de page
			if pageSize > MaxPageSize {
				pageSize = MaxPageSize
			}
		}
	}

	return page, pageSize
}

// getIDParam extrait un paramètre ID de l'URL
func getIDParam(r *http.Request, param string) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars[param], 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// getRequiredUserID récupère l'ID utilisateur du contexte ou renvoie une erreur
func getRequiredUserID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		respondWithError(w, http.StatusUnauthorized, "Utilisateur non authentifié")
		return 0, false
	}
	return userID, true
}

// extractFilters extrait les filtres de la requête
func extractFilters(r *http.Request, allowedFilters []string) map[string]string {
	filters := make(map[string]string)

	// Parcourir les paramètres de requête
	for key, values := range r.URL.Query() {
		// Vérifier si le filtre est autorisé
		for _, allowed := range allowedFilters {
			if strings.EqualFold(key, allowed) && len(values) > 0 {
				filters[allowed] = values[0]
				break
			}
		}
	}

	return filters
}
