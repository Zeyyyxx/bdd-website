package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter est un wrapper pour http.ResponseWriter qui capture le code de statut
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader capture le code de statut
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging est un middleware qui enregistre les requêtes HTTP
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Créer un wrapper pour capturer le code de statut
		rw := &responseWriter{w, http.StatusOK}

		// Enregistrer l'heure de début
		start := time.Now()

		// Traiter la requête
		next.ServeHTTP(rw, r)

		// Calculer la durée
		duration := time.Since(start)

		// Log de la requête
		log.Printf(
			"[%s] %s %s %d %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			rw.statusCode,
			duration,
		)
	})
}
