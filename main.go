package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"bdd-website/config"
	"bdd-website/internal/database"
	"bdd-website/internal/handlers"
	"bdd-website/internal/middleware"
)

func main() {
	// Charger la configuration
	cfg := config.LoadConfig()

	// Initialiser la base de données
	db, err := database.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de la base de données: %v", err)
	}
	defer db.Close()

	// Créer le routeur
	router := mux.NewRouter()

	// Middleware global
	router.Use(middleware.Logging)

	// Fichiers statiques
	fs := http.FileServer(http.Dir("./assets"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Routes de pages
	router.HandleFunc("/", handlers.HomePage).Methods("GET")
	router.HandleFunc("/about", handlers.AboutPage).Methods("GET")
	router.HandleFunc("/contact", handlers.ContactPage).Methods("GET")
	router.HandleFunc("/activities", handlers.ActivitiesPage).Methods("GET")
	router.HandleFunc("/login", handlers.LoginPage).Methods("GET")
	router.HandleFunc("/signup", handlers.SignupPage).Methods("GET")
	router.HandleFunc("/profile", handlers.ProfilePage).Methods("GET")

	// Routes d'authentification
	router.HandleFunc("/api/auth/register", handlers.Register(db)).Methods("POST")
	router.HandleFunc("/api/auth/login", handlers.Login(db, cfg.JWTSecret, cfg.JWTExpirationHours)).Methods("POST")

	// Routes utilisateurs
	userRouter := router.PathPrefix("/api/users").Subrouter()
	userRouter.Use(middleware.Auth(cfg.JWTSecret))
	userRouter.HandleFunc("/profile", handlers.GetUserProfile(db)).Methods("GET")
	userRouter.HandleFunc("/profile", handlers.UpdateUserProfile(db)).Methods("PUT")

	// Routes activités
	router.HandleFunc("/api/activities", handlers.GetActivities(db)).Methods("GET")
	router.HandleFunc("/api/activities/{id}", handlers.GetActivity(db)).Methods("GET")

	// Routes d'inscription aux activités (protégées)
	activityRegistrationRouter := router.PathPrefix("/api/activities").Subrouter()
	activityRegistrationRouter.Use(middleware.Auth(cfg.JWTSecret))
	activityRegistrationRouter.HandleFunc("/{id}/register", handlers.RegisterToActivity(db)).Methods("POST")
	activityRegistrationRouter.HandleFunc("/{id}/unregister", handlers.UnregisterFromActivity(db)).Methods("DELETE")

	// Routes contact
	router.HandleFunc("/api/contact", handlers.SubmitContactForm(db)).Methods("POST")

	// Routes du tableau de bord écologique (protégées)
	ecoDashboardRouter := router.PathPrefix("/api/eco-dashboard").Subrouter()
	ecoDashboardRouter.Use(middleware.Auth(cfg.JWTSecret))
	ecoDashboardRouter.HandleFunc("/points", handlers.GetUserEcoPoints(db)).Methods("GET")
	ecoDashboardRouter.HandleFunc("/challenges", handlers.GetUserChallenges(db)).Methods("GET")
	ecoDashboardRouter.HandleFunc("/challenges/{id}/join", handlers.JoinChallenge(db)).Methods("POST")
	ecoDashboardRouter.HandleFunc("/challenges/{id}/complete", handlers.CompleteChallenge(db)).Methods("POST")
	ecoDashboardRouter.HandleFunc("/badges", handlers.GetUserBadges(db)).Methods("GET")

	// Routes admin (protégées + vérification du rôle admin)
	adminRouter := router.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(middleware.Auth(cfg.JWTSecret))
	adminRouter.Use(middleware.AdminOnly)
	adminRouter.HandleFunc("/activities", handlers.AdminCreateActivity(db)).Methods("POST")
	adminRouter.HandleFunc("/activities/{id}", handlers.AdminUpdateActivity(db)).Methods("PUT")
	adminRouter.HandleFunc("/activities/{id}", handlers.AdminDeleteActivity(db)).Methods("DELETE")
	adminRouter.HandleFunc("/challenges", handlers.AdminCreateChallenge(db)).Methods("POST")
	adminRouter.HandleFunc("/challenges/{id}", handlers.AdminUpdateChallenge(db)).Methods("PUT")
	adminRouter.HandleFunc("/challenges/{id}", handlers.AdminDeleteChallenge(db)).Methods("DELETE")
	adminRouter.HandleFunc("/users", handlers.AdminGetUsers(db)).Methods("GET")
	adminRouter.HandleFunc("/contact-messages", handlers.AdminGetContactMessages(db)).Methods("GET")

	// Routes pages admin (protégées)
	adminPagesRouter := router.PathPrefix("/admin").Subrouter()
	adminPagesRouter.Use(middleware.Auth(cfg.JWTSecret))
	adminPagesRouter.Use(middleware.AdminOnly)
	adminPagesRouter.HandleFunc("", handlers.AdminDashboardPage).Methods("GET")
	adminPagesRouter.HandleFunc("/activities", handlers.AdminActivitiesPage).Methods("GET")
	adminPagesRouter.HandleFunc("/activities/edit/{id}", handlers.AdminEditActivityPage).Methods("GET")
	adminPagesRouter.HandleFunc("/activities/new", handlers.AdminNewActivityPage).Methods("GET")
	adminPagesRouter.HandleFunc("/challenges", handlers.AdminChallengesPage).Methods("GET")
	adminPagesRouter.HandleFunc("/users", handlers.AdminUsersPage).Methods("GET")
	adminPagesRouter.HandleFunc("/messages", handlers.AdminMessagesPage).Methods("GET")

	// Démarrer le serveur
	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Serveur démarré sur le port %d", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
