package database

import (
	"database/sql"
	"errors"
	"time"

	"bdd-website/internal/models"
	"bdd-website/internal/utils"
)

// CreateUser crée un nouvel utilisateur dans la base de données
func CreateUser(db *sql.DB, user models.UserRegister) (int64, error) {
	// Vérifier si l'email existe déjà
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email).Scan(&exists)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, errors.New("un utilisateur avec cet email existe déjà")
	}

	// Hacher le mot de passe
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	// Insérer l'utilisateur
	result, err := db.Exec(
		"INSERT INTO users (email, username, password_hash) VALUES (?, ?, ?)",
		user.Email, user.Username, hashedPassword,
	)
	if err != nil {
		return 0, err
	}

	// Récupérer l'ID généré
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Attribuer automatiquement le badge "Débutant écolo"
	_, err = db.Exec(
		"INSERT INTO user_badges (user_id, badge_id) VALUES (?, (SELECT id FROM badges WHERE name = 'Débutant écolo'))",
		userID,
	)
	if err != nil {
		// Ne pas échouer l'enregistrement si l'attribution du badge échoue
		// Juste enregistrer l'erreur (dans un vrai système, on utiliserait un logger)
	}

	return userID, nil
}

// GetUserByEmail récupère un utilisateur par son email
func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	user := &models.User{}

	err := db.QueryRow(
		"SELECT id, email, username, password_hash, is_admin, created_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.IsAdmin, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("utilisateur non trouvé")
		}
		return nil, err
	}

	return user, nil
}

// GetUserByID récupère un utilisateur par son ID
func GetUserByID(db *sql.DB, userID int64) (*models.User, error) {
	user := &models.User{}

	err := db.QueryRow(
		"SELECT id, email, username, password_hash, is_admin, created_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.IsAdmin, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("utilisateur non trouvé")
		}
		return nil, err
	}

	return user, nil
}

// GetUserProfile récupère le profil complet d'un utilisateur avec des statistiques
func GetUserProfile(db *sql.DB, userID int64) (*models.UserProfile, error) {
	// Récupérer les informations de base
	user, err := GetUserByID(db, userID)
	if err != nil {
		return nil, err
	}

	profile := &models.UserProfile{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
	}

	// Récupérer le nombre total de points écologiques
	err = db.QueryRow("SELECT COALESCE(SUM(points), 0) FROM eco_points WHERE user_id = ?", userID).Scan(&profile.TotalEcoPoints)
	if err != nil {
		// Ne pas échouer si cette requête échoue
		profile.TotalEcoPoints = 0
	}

	// Récupérer le nombre d'activités auxquelles l'utilisateur est inscrit
	err = db.QueryRow("SELECT COUNT(*) FROM registrations WHERE user_id = ?", userID).Scan(&profile.ActivityCount)
	if err != nil {
		// Ne pas échouer si cette requête échoue
		profile.ActivityCount = 0
	}

	// Récupérer le nombre de badges obtenus
	err = db.QueryRow("SELECT COUNT(*) FROM user_badges WHERE user_id = ?", userID).Scan(&profile.BadgeCount)
	if err != nil {
		// Ne pas échouer si cette requête échoue
		profile.BadgeCount = 0
	}

	return profile, nil
}

// UpdateUserProfile met à jour le profil d'un utilisateur
func UpdateUserProfile(db *sql.DB, userID int64, update models.UserProfileUpdate) error {
	// Vérifier si l'email est déjà utilisé par un autre utilisateur
	if update.Email != "" {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND id != ?)", update.Email, userID).Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			return errors.New("cet email est déjà utilisé par un autre utilisateur")
		}
	}

	// Commencer une transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Mise à jour des champs non-mot de passe
	if update.Username != "" || update.Email != "" {
		query := "UPDATE users SET"
		args := []interface{}{}

		if update.Username != "" && update.Email != "" {
			query += " username = ?, email = ?"
			args = append(args, update.Username, update.Email)
		} else if update.Username != "" {
			query += " username = ?"
			args = append(args, update.Username)
		} else if update.Email != "" {
			query += " email = ?"
			args = append(args, update.Email)
		}

		query += " WHERE id = ?"
		args = append(args, userID)

		_, err := tx.Exec(query, args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Mise à jour du mot de passe si fourni
	if update.Password != "" {
		hashedPassword, err := utils.HashPassword(update.Password)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec("UPDATE users SET password_hash = ? WHERE id = ?", hashedPassword, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit de la transaction
	return tx.Commit()
}

// GetAllUsers récupère tous les utilisateurs (pour l'admin)
func GetAllUsers(db *sql.DB, page, pageSize int) ([]models.UserProfile, int, error) {
	// Calculer l'offset pour la pagination
	offset := (page - 1) * pageSize

	// Récupérer le nombre total d'utilisateurs
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Récupérer les utilisateurs avec pagination
	rows, err := db.Query(`
		SELECT u.id, u.email, u.username, u.is_admin, u.created_at,
		       COALESCE((SELECT SUM(points) FROM eco_points WHERE user_id = u.id), 0) as total_points,
		       (SELECT COUNT(*) FROM registrations WHERE user_id = u.id) as activity_count,
		       (SELECT COUNT(*) FROM user_badges WHERE user_id = u.id) as badge_count
		FROM users u
		ORDER BY u.created_at DESC
		LIMIT ? OFFSET ?
	`, pageSize, offset)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := []models.UserProfile{}
	for rows.Next() {
		var user models.UserProfile
		var createdAt time.Time

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.IsAdmin,
			&createdAt,
			&user.TotalEcoPoints,
			&user.ActivityCount,
			&user.BadgeCount,
		)

		if err != nil {
			return nil, 0, err
		}

		user.CreatedAt = createdAt
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUserAdminStatus met à jour le statut d'administrateur d'un utilisateur
func UpdateUserAdminStatus(db *sql.DB, userID int64, isAdmin bool) error {
	_, err := db.Exec("UPDATE users SET is_admin = ? WHERE id = ?", isAdmin, userID)
	return err
}

// GetAdminStats récupère les statistiques pour le tableau de bord administrateur
func GetAdminStats(db *sql.DB) (*models.AdminStats, error) {
	stats := &models.AdminStats{}

	// Compter les utilisateurs
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.UsersCount)
	if err != nil {
		return nil, err
	}

	// Compter les activités
	err = db.QueryRow("SELECT COUNT(*) FROM activities").Scan(&stats.ActivitiesCount)
	if err != nil {
		return nil, err
	}

	// Compter les défis
	err = db.QueryRow("SELECT COUNT(*) FROM eco_challenges").Scan(&stats.ChallengesCount)
	if err != nil {
		return nil, err
	}

	// Compter les messages non lus
	err = db.QueryRow("SELECT COUNT(*) FROM contact_messages WHERE is_read = 0").Scan(&stats.UnreadMessagesCount)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
