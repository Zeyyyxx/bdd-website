package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

// serveTemplate sert un fichier de template HTML
func serveTemplate(w http.ResponseWriter, r *http.Request, filename string) {
	// Chemin complet du fichier
	path := filepath.Join("./templates", filename)

	// Vérifier si le fichier existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Si le fichier n'existe pas, renvoyer 404
		http.NotFound(w, r)
		return
	}

	// Servir le fichier HTML
	http.ServeFile(w, r, path)
}

// HomePage sert la page d'accueil
func HomePage(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, r, "index.html")
}

// AboutPage sert la page "Qui sommes-nous?"
func AboutPage(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, r, "about.html")
}

// ContactPage sert la page de contact
func ContactPage(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, r, "contact.html")
}

// ActivitiesPage sert la page des activités
func ActivitiesPage(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, r, "activities.html")
}

// LoginPage sert la page de connexion
func LoginPage(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, r, "login.html")
}

// SignupPage sert la page d'inscription
func SignupPage(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, r, "signup.html")
}

// ProfilePage sert la page de profil utilisateur
func ProfilePage(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, r, "profile.html")
}

// AdminDashboardPage sert la page de tableau de bord admin
func AdminDashboardPage(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, r, "admin/dashboard.html")
}

// AdminActivitiesPage sert la page de gestion des activités admin
func AdminActivitiesPage(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, r, "admin/activities.html")
}

// AdminEditActivityPage sert la page d'édition d'activité admin
func AdminEditActivityPage(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, r, "admin/edit-activity.html")
}

// AdminNewActivityPage sert la page de création d'activité admin
func AdminNewActivityPage(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, r, "admin/new-activity.html")
}

// AdminChallengesPage sert la page de gestion des défis admin
func AdminChallengesPage(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, r, "admin/challenges.html")
}

// AdminUsersPage sert la page de gestion des utilisateurs admin
func AdminUsersPage(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, r, "admin/users.html")
}

// AdminMessagesPage sert la page de gestion des messages admin
func AdminMessagesPage(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, r, "admin/messages.html")
}
