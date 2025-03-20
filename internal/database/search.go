package database

import (
	"database/sql"
	"strings"

	"bdd-website/internal/models"
)

// SearchUsers recherche des utilisateurs basé sur un terme de recherche
func SearchUsers(db *sql.DB, query string, limit int) ([]models.UserProfile, error) {
	// Préparer le terme de recherche
	searchTerm := "%" + strings.ToLower(query) + "%"

	// Effectuer la recherche
	rows, err := db.Query(`
		SELECT u.id, u.email, u.username, u.is_admin, u.created_at,
		       COALESCE((SELECT SUM(points) FROM eco_points WHERE user_id = u.id), 0) as total_points,
		       (SELECT COUNT(*) FROM registrations WHERE user_id = u.id) as activity_count,
		       (SELECT COUNT(*) FROM user_badges WHERE user_id = u.id) as badge_count
		FROM users u
		WHERE LOWER(u.username) LIKE ? OR LOWER(u.email) LIKE ?
		ORDER BY 
			CASE 
				WHEN LOWER(u.username) = LOWER(?) THEN 0
				WHEN LOWER(u.email) = LOWER(?) THEN 0
				WHEN LOWER(u.username) LIKE ? THEN 1
				WHEN LOWER(u.email) LIKE ? THEN 1
				ELSE 2
			END,
			u.created_at DESC
		LIMIT ?
	`, searchTerm, searchTerm, query, query, searchTerm, searchTerm, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.UserProfile{}
	for rows.Next() {
		var user models.UserProfile
		var createdAt sql.NullTime

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
			return nil, err
		}

		if createdAt.Valid {
			user.CreatedAt = createdAt.Time
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
