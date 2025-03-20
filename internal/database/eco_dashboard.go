package database

import (
	"database/sql"
	"errors"
	"time"

	"bdd-website/internal/models"
)

// GetUserEcoPoints récupère les points écologiques d'un utilisateur
func GetUserEcoPoints(db *sql.DB, userID int64) ([]models.EcoPoint, int, error) {
	// Récupérer les points
	rows, err := db.Query(`
		SELECT ep.id, ep.user_id, ep.activity_id, ep.challenge_id, 
		       ep.points, ep.description, ep.date,
		       a.title as activity_title, c.title as challenge_title
		FROM eco_points ep
		LEFT JOIN activities a ON ep.activity_id = a.id
		LEFT JOIN eco_challenges c ON ep.challenge_id = c.id
		WHERE ep.user_id = ?
		ORDER BY ep.date DESC
	`, userID)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Parcourir les résultats
	points := []models.EcoPoint{}
	for rows.Next() {
		var point models.EcoPoint
		var date time.Time
		var activityID, challengeID sql.NullInt64
		var activityTitle, challengeTitle sql.NullString

		err := rows.Scan(
			&point.ID, &point.UserID, &activityID, &challengeID,
			&point.Points, &point.Description, &date,
			&activityTitle, &challengeTitle,
		)

		if err != nil {
			return nil, 0, err
		}

		// Convertir les valeurs nullables
		if activityID.Valid {
			point.ActivityID = activityID.Int64
		}

		if challengeID.Valid {
			point.ChallengeID = challengeID.Int64
		}

		if activityTitle.Valid {
			point.ActivityTitle = activityTitle.String
		}

		if challengeTitle.Valid {
			point.ChallengeTitle = challengeTitle.String
		}

		point.Date = date
		points = append(points, point)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Calculer le total des points
	var totalPoints int
	err = db.QueryRow("SELECT COALESCE(SUM(points), 0) FROM eco_points WHERE user_id = ?", userID).Scan(&totalPoints)
	if err != nil {
		return nil, 0, err
	}

	return points, totalPoints, nil
}

// AddEcoPoints ajoute des points écologiques à un utilisateur
func AddEcoPoints(db *sql.DB, userID int64, activityID, challengeID int64, points int, description string) (int64, error) {
	// Vérifier que les points sont positifs
	if points <= 0 {
		return 0, errors.New("les points doivent être positifs")
	}

	// Insérer les points
	result, err := db.Exec(
		"INSERT INTO eco_points (user_id, activity_id, challenge_id, points, description) VALUES (?, ?, ?, ?, ?)",
		userID, nullIfZero(activityID), nullIfZero(challengeID), points, description,
	)

	if err != nil {
		return 0, err
	}

	// Récupérer l'ID généré
	pointID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Vérifier et attribuer les badges basés sur les points totaux
	go checkAndAwardBadges(db, userID)

	return pointID, nil
}

// nullIfZero retourne NULL si la valeur est 0
func nullIfZero(val int64) interface{} {
	if val == 0 {
		return nil
	}
	return val
}

// GetChallenges récupère les défis écologiques disponibles
func GetChallenges(db *sql.DB, userID int64, activeOnly bool) ([]models.Challenge, error) {
	// Construire la requête
	query := `
		SELECT c.id, c.title, c.description, c.points, c.duration_days, 
		       c.start_date, c.end_date, c.is_active, c.created_at,
		       cp.status, cp.joined_at, cp.completed_at
		FROM eco_challenges c
		LEFT JOIN challenge_participants cp ON c.id = cp.challenge_id AND cp.user_id = ?
	`

	// Ajouter les filtres
	whereClause := ""
	if activeOnly {
		whereClause = " WHERE c.is_active = 1"
	}

	// Ordonner
	orderClause := " ORDER BY c.created_at DESC"

	// Requête finale
	query = query + whereClause + orderClause

	// Exécuter la requête
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parcourir les résultats
	challenges := []models.Challenge{}
	for rows.Next() {
		var challenge models.Challenge
		var startDate, endDate, createdAt sql.NullTime
		var status sql.NullString
		var joinedAt, completedAt sql.NullTime

		err := rows.Scan(
			&challenge.ID, &challenge.Title, &challenge.Description, &challenge.Points,
			&challenge.DurationDays, &startDate, &endDate, &challenge.IsActive, &createdAt,
			&status, &joinedAt, &completedAt,
		)

		if err != nil {
			return nil, err
		}

		// Convertir les valeurs nullables
		if startDate.Valid {
			challenge.StartDate = startDate.Time
		}

		if endDate.Valid {
			challenge.EndDate = endDate.Time
		}

		if createdAt.Valid {
			challenge.CreatedAt = createdAt.Time
		}

		if status.Valid {
			challenge.UserStatus = status.String
		} else {
			challenge.UserStatus = "not_joined"
		}

		if joinedAt.Valid {
			challenge.JoinedAt = joinedAt.Time
		}

		if completedAt.Valid {
			challenge.CompletedAt = completedAt.Time
		}

		challenges = append(challenges, challenge)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return challenges, nil
}

// GetUserChallenges récupère les défis auxquels un utilisateur participe
func GetUserChallenges(db *sql.DB, userID int64) ([]models.Challenge, error) {
	// Exécuter la requête
	rows, err := db.Query(`
		SELECT c.id, c.title, c.description, c.points, c.duration_days, 
		       c.start_date, c.end_date, c.is_active, c.created_at,
		       cp.status, cp.joined_at, cp.completed_at
		FROM eco_challenges c
		JOIN challenge_participants cp ON c.id = cp.challenge_id AND cp.user_id = ?
		ORDER BY cp.joined_at DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parcourir les résultats
	challenges := []models.Challenge{}
	for rows.Next() {
		var challenge models.Challenge
		var startDate, endDate, createdAt sql.NullTime
		var status string
		var joinedAt time.Time
		var completedAt sql.NullTime

		err := rows.Scan(
			&challenge.ID, &challenge.Title, &challenge.Description, &challenge.Points,
			&challenge.DurationDays, &startDate, &endDate, &challenge.IsActive, &createdAt,
			&status, &joinedAt, &completedAt,
		)

		if err != nil {
			return nil, err
		}

		// Convertir les valeurs nullables
		if startDate.Valid {
			challenge.StartDate = startDate.Time
		}

		if endDate.Valid {
			challenge.EndDate = endDate.Time
		}

		if createdAt.Valid {
			challenge.CreatedAt = createdAt.Time
		}

		challenge.UserStatus = status
		challenge.JoinedAt = joinedAt

		if completedAt.Valid {
			challenge.CompletedAt = completedAt.Time
		}

		challenges = append(challenges, challenge)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return challenges, nil
}

// CreateChallenge crée un nouveau défi écologique
func CreateChallenge(db *sql.DB, challenge models.ChallengeCreate) (int64, error) {
	// Préparer les valeurs nullables
	var startDateArg, endDateArg interface{}

	if !challenge.StartDate.IsZero() {
		startDateArg = challenge.StartDate
	} else {
		startDateArg = nil
	}

	if !challenge.EndDate.IsZero() {
		endDateArg = challenge.EndDate
	} else {
		endDateArg = nil
	}

	// Insérer le défi
	result, err := db.Exec(
		`INSERT INTO eco_challenges 
		(title, description, points, duration_days, start_date, end_date, is_active) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		challenge.Title, challenge.Description, challenge.Points,
		challenge.DurationDays, startDateArg, endDateArg, challenge.IsActive,
	)

	if err != nil {
		return 0, err
	}

	// Récupérer l'ID généré
	return result.LastInsertId()
}

// UpdateChallenge met à jour un défi écologique
func UpdateChallenge(db *sql.DB, challengeID int64, challenge models.ChallengeUpdate) error {
	// Vérifier si le défi existe
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM eco_challenges WHERE id = ?)", challengeID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("défi non trouvé")
	}

	// Préparer les valeurs nullables
	var startDateArg, endDateArg interface{}

	if !challenge.StartDate.IsZero() {
		startDateArg = challenge.StartDate
	} else {
		startDateArg = nil
	}

	if !challenge.EndDate.IsZero() {
		endDateArg = challenge.EndDate
	} else {
		endDateArg = nil
	}

	// Mettre à jour le défi
	_, err = db.Exec(
		`UPDATE eco_challenges 
		SET title = ?, description = ?, points = ?, 
		    duration_days = ?, start_date = ?, end_date = ?, is_active = ?
		WHERE id = ?`,
		challenge.Title, challenge.Description, challenge.Points,
		challenge.DurationDays, startDateArg, endDateArg, challenge.IsActive,
		challengeID,
	)

	return err
}

// DeleteChallenge supprime un défi écologique
func DeleteChallenge(db *sql.DB, challengeID int64) error {
	// Vérifier si le défi existe
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM eco_challenges WHERE id = ?)", challengeID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("défi non trouvé")
	}

	// Supprimer dans une transaction pour gérer les dépendances
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Supprimer les participations liées
	_, err = tx.Exec("DELETE FROM challenge_participants WHERE challenge_id = ?", challengeID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Mettre à null les références dans eco_points
	_, err = tx.Exec("UPDATE eco_points SET challenge_id = NULL WHERE challenge_id = ?", challengeID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Supprimer le défi
	_, err = tx.Exec("DELETE FROM eco_challenges WHERE id = ?", challengeID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// JoinChallenge permet à un utilisateur de rejoindre un défi
func JoinChallenge(db *sql.DB, userID, challengeID int64) error {
	// Vérifier si le défi existe et est actif
	var isActive bool
	err := db.QueryRow("SELECT is_active FROM eco_challenges WHERE id = ?", challengeID).Scan(&isActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("défi non trouvé")
		}
		return err
	}

	if !isActive {
		return errors.New("ce défi n'est pas actif")
	}

	// Vérifier si l'utilisateur participe déjà
	var status string
	err = db.QueryRow(
		"SELECT status FROM challenge_participants WHERE user_id = ? AND challenge_id = ?",
		userID, challengeID,
	).Scan(&status)

	if err == nil {
		// L'utilisateur participe déjà
		if status == "in_progress" || status == "completed" {
			return errors.New("vous participez déjà à ce défi")
		}

		// Si abandonné, mettre à jour le statut
		_, err = db.Exec(
			"UPDATE challenge_participants SET status = 'in_progress', joined_at = ?, completed_at = NULL WHERE user_id = ? AND challenge_id = ?",
			time.Now(), userID, challengeID,
		)
		return err
	} else if err != sql.ErrNoRows {
		return err
	}

	// Sinon, créer une nouvelle participation
	_, err = db.Exec(
		"INSERT INTO challenge_participants (user_id, challenge_id, status, joined_at) VALUES (?, ?, 'in_progress', ?)",
		userID, challengeID, time.Now(),
	)

	return err
}

// CompleteChallenge marque un défi comme terminé pour un utilisateur
func CompleteChallenge(db *sql.DB, userID, challengeID int64) error {
	// Vérifier si l'utilisateur participe au défi
	var status string
	var joinedAt time.Time
	var points int

	err := db.QueryRow(`
		SELECT cp.status, cp.joined_at, c.points
		FROM challenge_participants cp
		JOIN eco_challenges c ON cp.challenge_id = c.id
		WHERE cp.user_id = ? AND cp.challenge_id = ?
	`, userID, challengeID).Scan(&status, &joinedAt, &points)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("vous ne participez pas à ce défi")
		}
		return err
	}

	if status != "in_progress" {
		if status == "completed" {
			return errors.New("vous avez déjà terminé ce défi")
		}
		return errors.New("vous ne pouvez pas terminer ce défi")
	}

	// Commencer une transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Mettre à jour le statut
	now := time.Now()
	_, err = tx.Exec(
		"UPDATE challenge_participants SET status = 'completed', completed_at = ? WHERE user_id = ? AND challenge_id = ?",
		now, userID, challengeID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Ajouter les points
	_, err = tx.Exec(
		"INSERT INTO eco_points (user_id, challenge_id, points, description) VALUES (?, ?, ?, ?)",
		userID, challengeID, points, "Défi complété",
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit de la transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	// Vérifier et attribuer les badges basés sur les points totaux
	go checkAndAwardBadges(db, userID)

	return nil
}

// GetUserBadges récupère les badges d'un utilisateur
func GetUserBadges(db *sql.DB, userID int64) ([]models.Badge, []models.Badge, error) {
	// Récupérer les points totaux de l'utilisateur
	var totalPoints int
	err := db.QueryRow("SELECT COALESCE(SUM(points), 0) FROM eco_points WHERE user_id = ?", userID).Scan(&totalPoints)
	if err != nil {
		return nil, nil, err
	}

	// Récupérer les badges déjà obtenus
	earnedBadgesRows, err := db.Query(`
		SELECT b.id, b.name, b.description, b.image_path, b.required_points, b.category, ub.earned_at
		FROM badges b
		JOIN user_badges ub ON b.id = ub.badge_id
		WHERE ub.user_id = ?
		ORDER BY ub.earned_at DESC
	`, userID)

	if err != nil {
		return nil, nil, err
	}
	defer earnedBadgesRows.Close()

	// Parcourir les badges obtenus
	earnedBadges := []models.Badge{}
	for earnedBadgesRows.Next() {
		var badge models.Badge
		var earnedAt time.Time

		err := earnedBadgesRows.Scan(
			&badge.ID, &badge.Name, &badge.Description, &badge.ImagePath,
			&badge.RequiredPoints, &badge.Category, &earnedAt,
		)

		if err != nil {
			return nil, nil, err
		}

		badge.IsEarned = true
		badge.EarnedAt = earnedAt
		earnedBadges = append(earnedBadges, badge)
	}

	if err = earnedBadgesRows.Err(); err != nil {
		return nil, nil, err
	}

	// Récupérer les badges disponibles mais non obtenus
	availableBadgesRows, err := db.Query(`
		SELECT b.id, b.name, b.description, b.image_path, b.required_points, b.category
		FROM badges b
		WHERE b.id NOT IN (SELECT badge_id FROM user_badges WHERE user_id = ?)
		ORDER BY b.required_points ASC
	`, userID)

	if err != nil {
		return nil, nil, err
	}
	defer availableBadgesRows.Close()

	// Parcourir les badges disponibles
	availableBadges := []models.Badge{}
	for availableBadgesRows.Next() {
		var badge models.Badge

		err := availableBadgesRows.Scan(
			&badge.ID, &badge.Name, &badge.Description, &badge.ImagePath,
			&badge.RequiredPoints, &badge.Category,
		)

		if err != nil {
			return nil, nil, err
		}

		badge.IsEarned = false
		badge.PointsToGo = badge.RequiredPoints - totalPoints
		if badge.PointsToGo < 0 {
			badge.PointsToGo = 0
		}

		availableBadges = append(availableBadges, badge)
	}

	if err = availableBadgesRows.Err(); err != nil {
		return nil, nil, err
	}

	return earnedBadges, availableBadges, nil
}

// checkAndAwardBadges vérifie et attribue les badges en fonction des points
func checkAndAwardBadges(db *sql.DB, userID int64) {
	// Récupérer les points totaux de l'utilisateur
	var totalPoints int
	err := db.QueryRow("SELECT COALESCE(SUM(points), 0) FROM eco_points WHERE user_id = ?", userID).Scan(&totalPoints)
	if err != nil {
		return
	}

	// Récupérer les badges déjà obtenus
	rows, err := db.Query("SELECT badge_id FROM user_badges WHERE user_id = ?", userID)
	if err != nil {
		return
	}

	earnedBadgeIDs := make(map[int64]bool)
	for rows.Next() {
		var badgeID int64
		if err := rows.Scan(&badgeID); err != nil {
			rows.Close()
			return
		}
		earnedBadgeIDs[badgeID] = true
	}
	rows.Close()

	// Récupérer les badges disponibles
	badgeRows, err := db.Query(
		"SELECT id FROM badges WHERE required_points <= ? ORDER BY required_points ASC",
		totalPoints,
	)
	if err != nil {
		return
	}
	defer badgeRows.Close()

	// Attribuer les nouveaux badges
	for badgeRows.Next() {
		var badgeID int64
		if err := badgeRows.Scan(&badgeID); err != nil {
			return
		}

		// Si le badge n'a pas déjà été attribué
		if !earnedBadgeIDs[badgeID] {
			_, err := db.Exec(
				"INSERT INTO user_badges (user_id, badge_id) VALUES (?, ?)",
				userID, badgeID,
			)
			if err != nil {
				// Ignorer les erreurs d'insertion (comme les tentatives en double)
				continue
			}
		}
	}
}

// GetEcoDashboardSummary récupère un résumé du tableau de bord écologique
func GetEcoDashboardSummary(db *sql.DB, userID int64) (*models.EcoDashboardSummary, error) {
	summary := &models.EcoDashboardSummary{}

	// Récupérer les points totaux
	err := db.QueryRow("SELECT COALESCE(SUM(points), 0) FROM eco_points WHERE user_id = ?", userID).Scan(&summary.TotalPoints)
	if err != nil {
		return nil, err
	}

	// Récupérer le nombre d'activités auxquelles l'utilisateur a participé
	err = db.QueryRow("SELECT COUNT(DISTINCT activity_id) FROM registrations WHERE user_id = ?", userID).Scan(&summary.ActivitiesAttended)
	if err != nil {
		return nil, err
	}

	// Récupérer le nombre de défis complétés
	err = db.QueryRow("SELECT COUNT(*) FROM challenge_participants WHERE user_id = ? AND status = 'completed'", userID).Scan(&summary.ChallengesCompleted)
	if err != nil {
		return nil, err
	}

	// Récupérer le nombre de badges obtenus
	err = db.QueryRow("SELECT COUNT(*) FROM user_badges WHERE user_id = ?", userID).Scan(&summary.BadgesEarned)
	if err != nil {
		return nil, err
	}

	// Récupérer le classement de l'utilisateur et le nombre total d'utilisateurs
	err = db.QueryRow(`
		SELECT ranking, total_users
		FROM (
			SELECT 
				user_id, 
				COALESCE(SUM(points), 0) as total_points,
				RANK() OVER (ORDER BY COALESCE(SUM(points), 0) DESC) as ranking,
				COUNT(*) OVER () as total_users
			FROM eco_points
			GROUP BY user_id
		) rankings
		WHERE user_id = ?
	`, userID).Scan(&summary.Ranking, &summary.TotalUsers)

	if err == sql.ErrNoRows {
		// L'utilisateur n'a pas encore de points, donc dernier du classement
		var totalUsers int
		err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
		if err != nil {
			return nil, err
		}

		summary.Ranking = totalUsers
		summary.TotalUsers = totalUsers
	} else if err != nil {
		return nil, err
	}

	return summary, nil
}
