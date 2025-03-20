package models

import (
	"time"
)

// User représente un utilisateur du système
type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Jamais envoyé dans les réponses JSON
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRegister représente les données requises pour l'inscription d'un utilisateur
type UserRegister struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserLogin représente les données requises pour la connexion d'un utilisateur
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserProfile représente les données du profil utilisateur exposées à l'API
type UserProfile struct {
	ID             int64     `json:"id"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	IsAdmin        bool      `json:"is_admin"`
	CreatedAt      time.Time `json:"created_at"`
	TotalEcoPoints int       `json:"total_eco_points"`
	ActivityCount  int       `json:"activity_count"`
	BadgeCount     int       `json:"badge_count"`
}

// UserProfileUpdate représente les données modifiables du profil utilisateur
type UserProfileUpdate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // Optionnel
}

// UserResponse représente la réponse après authentification
type UserResponse struct {
	User  UserProfile `json:"user"`
	Token string      `json:"token"`
}

// Activity représente une activité ou un événement du BDD
type Activity struct {
	ID                  int64     `json:"id"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	ImagePath           string    `json:"image_path"`
	StartDate           time.Time `json:"start_date"`
	EndDate             time.Time `json:"end_date"`
	Location            string    `json:"location"`
	MaxParticipants     int       `json:"max_participants"`
	EcoPoints           int       `json:"eco_points"`
	CurrentParticipants int       `json:"current_participants,omitempty"`
	UserRegistered      bool      `json:"user_registered,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ActivityCreate représente les données pour créer une nouvelle activité
type ActivityCreate struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	ImagePath       string    `json:"image_path"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Location        string    `json:"location"`
	MaxParticipants int       `json:"max_participants"`
	EcoPoints       int       `json:"eco_points"`
}

// ActivityUpdate représente les données pour mettre à jour une activité
type ActivityUpdate struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	ImagePath       string    `json:"image_path"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Location        string    `json:"location"`
	MaxParticipants int       `json:"max_participants"`
	EcoPoints       int       `json:"eco_points"`
}

// ActivitiesResponse représente la réponse de la liste des activités
type ActivitiesResponse struct {
	Activities []Activity `json:"activities"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
}

// ContactMessage représente un message envoyé via le formulaire de contact
type ContactMessage struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Subject     string    `json:"subject"`
	Message     string    `json:"message"`
	SubmittedAt time.Time `json:"submitted_at"`
	IsRead      bool      `json:"is_read"`
}

// ContactMessageCreate représente les données pour créer un nouveau message de contact
type ContactMessageCreate struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// ContactMessagesResponse représente la réponse pour les messages de contact
type ContactMessagesResponse struct {
	Messages []ContactMessage `json:"messages"`
	Total    int              `json:"total"`
	Unread   int              `json:"unread"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

// EcoPoint représente un point d'impact écologique gagné par un utilisateur
type EcoPoint struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	ActivityID     int64     `json:"activity_id,omitempty"`
	ChallengeID    int64     `json:"challenge_id,omitempty"`
	Points         int       `json:"points"`
	Description    string    `json:"description"`
	Date           time.Time `json:"date"`
	ActivityTitle  string    `json:"activity_title,omitempty"`
	ChallengeTitle string    `json:"challenge_title,omitempty"`
}

// EcoPointsResponse représente la réponse pour les points écologiques d'un utilisateur
type EcoPointsResponse struct {
	Points      []EcoPoint `json:"points"`
	TotalPoints int        `json:"total_points"`
}

// Challenge représente un défi écologique
type Challenge struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Points       int       `json:"points"`
	DurationDays int       `json:"duration_days"`
	StartDate    time.Time `json:"start_date,omitempty"`
	EndDate      time.Time `json:"end_date,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UserStatus   string    `json:"user_status,omitempty"` // 'not_joined', 'in_progress', 'completed', 'abandoned'
	JoinedAt     time.Time `json:"joined_at,omitempty"`
	CompletedAt  time.Time `json:"completed_at,omitempty"`
}

// ChallengeCreate représente les données pour créer un nouveau défi
type ChallengeCreate struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Points       int       `json:"points"`
	DurationDays int       `json:"duration_days"`
	StartDate    time.Time `json:"start_date,omitempty"`
	EndDate      time.Time `json:"end_date,omitempty"`
	IsActive     bool      `json:"is_active"`
}

// ChallengeUpdate représente les données pour mettre à jour un défi
type ChallengeUpdate struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Points       int       `json:"points"`
	DurationDays int       `json:"duration_days"`
	StartDate    time.Time `json:"start_date,omitempty"`
	EndDate      time.Time `json:"end_date,omitempty"`
	IsActive     bool      `json:"is_active"`
}

// ChallengesResponse représente la réponse pour les défis écologiques
type ChallengesResponse struct {
	ActiveChallenges    []Challenge `json:"active_challenges"`
	UserChallenges      []Challenge `json:"user_challenges"`
	CompletedChallenges []Challenge `json:"completed_challenges"`
}

// Badge représente un badge écologique
type Badge struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ImagePath      string    `json:"image_path"`
	RequiredPoints int       `json:"required_points"`
	Category       string    `json:"category"`
	IsEarned       bool      `json:"is_earned,omitempty"`
	EarnedAt       time.Time `json:"earned_at,omitempty"`
	PointsToGo     int       `json:"points_to_go,omitempty"` // Points manquants pour obtenir le badge
}

// BadgesResponse représente la réponse pour les badges écologiques
type BadgesResponse struct {
	EarnedBadges    []Badge `json:"earned_badges"`
	AvailableBadges []Badge `json:"available_badges"`
}

// EcoDashboardSummary représente un résumé du tableau de bord écologique
type EcoDashboardSummary struct {
	TotalPoints         int `json:"total_points"`
	ActivitiesAttended  int `json:"activities_attended"`
	ChallengesCompleted int `json:"challenges_completed"`
	BadgesEarned        int `json:"badges_earned"`
	Ranking             int `json:"ranking,omitempty"`     // Position dans le classement général
	TotalUsers          int `json:"total_users,omitempty"` // Nombre total d'utilisateurs pour le classement
}

// AdminStats représente les statistiques pour le tableau de bord administrateur
type AdminStats struct {
	UsersCount          int `json:"users_count"`
	ActivitiesCount     int `json:"activities_count"`
	ChallengesCount     int `json:"challenges_count"`
	UnreadMessagesCount int `json:"unread_messages_count"`
}
