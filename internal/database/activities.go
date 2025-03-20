package database

import (
	"database/sql"
	"errors"
	"time"

	"bdd-website/internal/models"
)

// CreateActivity crée une nouvelle activité dans la base de données
func CreateActivity(db *sql.DB, activity models.ActivityCreate) (int64, error) {
	// Insérer l'activité
	result, err := db.Exec(
		`INSERT INTO activities 
		(title, description, image_path, start_date, end_date, location, max_participants, eco_points, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		activity.Title, activity.Description, activity.ImagePath,
		activity.StartDate, activity.EndDate, activity.Location,
		activity.MaxParticipants, activity.EcoPoints, time.Now(),
	)

	if err != nil {
		return 0, err
	}

	// Récupérer l'ID généré
	return result.LastInsertId()
}

// UpdateActivity met à jour une activité existante
func UpdateActivity(db *sql.DB, activityID int64, activity models.ActivityUpdate) error {
	// Vérifier si l'activité existe
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM activities WHERE id = ?)", activityID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("activité non trouvée")
	}

	// Mettre à jour l'activité
	_, err = db.Exec(
		`UPDATE activities 
		SET title = ?, description = ?, image_path = ?, 
		    start_date = ?, end_date = ?, location = ?, 
		    max_participants = ?, eco_points = ?, updated_at = ?
		WHERE id = ?`,
		activity.Title, activity.Description, activity.ImagePath,
		activity.StartDate, activity.EndDate, activity.Location,
		activity.MaxParticipants, activity.EcoPoints, time.Now(),
		activityID,
	)

	return err
}

// DeleteActivity supprime une activité
func DeleteActivity(db *sql.DB, activityID int64) error {
	// Vérifier si l'activité existe
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM activities WHERE id = ?)", activityID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("activité non trouvée")
	}

	// Supprimer dans une transaction pour gérer les dépendances
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Supprimer les inscriptions liées
	_, err = tx.Exec("DELETE FROM registrations WHERE activity_id = ?", activityID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Mettre à null les références dans eco_points
	_, err = tx.Exec("UPDATE eco_points SET activity_id = NULL WHERE activity_id = ?", activityID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Supprimer l'activité
	_, err = tx.Exec("DELETE FROM activities WHERE id = ?", activityID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetActivities récupère les activités avec pagination et filtres
func GetActivities(db *sql.DB, page, pageSize int, upcoming bool, userID int64) ([]models.Activity, int, error) {
	// Calculer l'offset pour la pagination
	offset := (page - 1) * pageSize

	// Base de la requête
	baseQuery := `
		SELECT a.id, a.title, a.description, a.image_path, 
		       a.start_date, a.end_date, a.location, 
		       a.max_participants, a.eco_points, a.created_at, a.updated_at,
		       COUNT(r.id) as current_participants
	`

	// Ajouter un champ pour indiquer si l'utilisateur est inscrit (si userID > 0)
	if userID > 0 {
		baseQuery += `, 
		(SELECT EXISTS(SELECT 1 FROM registrations WHERE user_id = ? AND activity_id = a.id)) as user_registered
		`
	}

	// Compléter la requête
	baseQuery += `
		FROM activities a
		LEFT JOIN registrations r ON a.id = r.activity_id
	`

	// Ajouter les filtres
	whereClause := ""
	if upcoming {
		whereClause = " WHERE a.end_date >= ?"
	}

	// Grouper et ordonner
	groupAndOrder := `
		GROUP BY a.id
		ORDER BY a.start_date ASC
		LIMIT ? OFFSET ?
	`

	// Construire la requête finale
	query := baseQuery + whereClause + groupAndOrder

	// Préparer les arguments
	var args []interface{}
	if userID > 0 {
		args = append(args, userID)
	}
	if upcoming {
		args = append(args, time.Now())
	}
	args = append(args, pageSize, offset)

	// Exécuter la requête
	var rows *sql.Rows
	var err error
	rows, err = db.Query(query, args...)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Parcourir les résultats
	activities := []models.Activity{}
	for rows.Next() {
		var activity models.Activity
		var startDate, endDate, createdAt, updatedAt time.Time
		var userRegistered sql.NullBool

		// Préparer les variables pour le scan
		scanArgs := []interface{}{
			&activity.ID, &activity.Title, &activity.Description, &activity.ImagePath,
			&startDate, &endDate, &activity.Location,
			&activity.MaxParticipants, &activity.EcoPoints, &createdAt, &updatedAt,
			&activity.CurrentParticipants,
		}

		// Ajouter userRegistered si nécessaire
		if userID > 0 {
			scanArgs = append(scanArgs, &userRegistered)
		}

		// Scanner les résultats
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, 0, err
		}

		// Convertir les dates
		activity.StartDate = startDate
		activity.EndDate = endDate
		activity.CreatedAt = createdAt
		activity.UpdatedAt = updatedAt

		// Assigner userRegistered si nécessaire
		if userID > 0 && userRegistered.Valid {
			activity.UserRegistered = userRegistered.Bool
		}

		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Récupérer le total pour la pagination
	var total int
	countQuery := "SELECT COUNT(*) FROM activities"
	if upcoming {
		countQuery += " WHERE end_date >= ?"
		err = db.QueryRow(countQuery, time.Now()).Scan(&total)
	} else {
		err = db.QueryRow(countQuery).Scan(&total)
	}

	if err != nil {
		return nil, 0, err
	}

	return activities, total, nil
}

// GetActivity récupère les détails d'une activité spécifique
func GetActivity(db *sql.DB, activityID int64, userID int64) (*models.Activity, error) {
	// Construire la requête
	query := `
		SELECT a.id, a.title, a.description, a.image_path, 
		       a.start_date, a.end_date, a.location, 
		       a.max_participants, a.eco_points, a.created_at, a.updated_at,
		       COUNT(r.id) as current_participants
	`

	// Ajouter un champ pour indiquer si l'utilisateur est inscrit (si userID > 0)
	if userID > 0 {
		query += `, 
		(SELECT EXISTS(SELECT 1 FROM registrations WHERE user_id = ? AND activity_id = a.id)) as user_registered
		`
	}

	// Compléter la requête
	query += `
		FROM activities a
		LEFT JOIN registrations r ON a.id = r.activity_id
		WHERE a.id = ?
		GROUP BY a.id
	`

	// Préparer les arguments
	var args []interface{}
	if userID > 0 {
		args = append(args, userID)
	}
	args = append(args, activityID)

	// Exécuter la requête
	var activity models.Activity
	var startDate, endDate, createdAt, updatedAt time.Time
	var userRegistered sql.NullBool

	// Préparer les variables pour le scan
	scanArgs := []interface{}{
		&activity.ID, &activity.Title, &activity.Description, &activity.ImagePath,
		&startDate, &endDate, &activity.Location,
		&activity.MaxParticipants, &activity.EcoPoints, &createdAt, &updatedAt,
		&activity.CurrentParticipants,
	}

	// Ajouter userRegistered si nécessaire
	if userID > 0 {
		scanArgs = append(scanArgs, &userRegistered)
	}

	// Exécuter la requête
	err := db.QueryRow(query, args...).Scan(scanArgs...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("activité non trouvée")
		}
		return nil, err
	}

	// Convertir les dates
	activity.StartDate = startDate
	activity.EndDate = endDate
	activity.CreatedAt = createdAt
	activity.UpdatedAt = updatedAt

	// Assigner userRegistered si nécessaire
	if userID > 0 && userRegistered.Valid {
		activity.UserRegistered = userRegistered.Bool
	}

	return &activity, nil
}

// RegisterToActivity inscrit un utilisateur à une activité
func RegisterToActivity(db *sql.DB, userID, activityID int64) error {
	// Vérifier si l'utilisateur est déjà inscrit
	var isRegistered bool
	err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM registrations WHERE user_id = ? AND activity_id = ?)",
		userID, activityID,
	).Scan(&isRegistered)

	if err != nil {
		return err
	}

	if isRegistered {
		return errors.New("vous êtes déjà inscrit à cette activité")
	}

	// Vérifier si l'activité existe et n'est pas déjà pleine
	var maxParticipants, currentParticipants int
	var activityStartDate time.Time
	var ecoPoints int

	err = db.QueryRow(`
		SELECT a.max_participants, COUNT(r.id), a.start_date, a.eco_points
		FROM activities a
		LEFT JOIN registrations r ON a.id = r.activity_id
		WHERE a.id = ?
		GROUP BY a.id
	`, activityID).Scan(&maxParticipants, &currentParticipants, &activityStartDate, &ecoPoints)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("activité non trouvée")
		}
		return err
	}

	if maxParticipants > 0 && currentParticipants >= maxParticipants {
		return errors.New("cette activité est complète")
	}

	// Vérifier si l'activité n'est pas déjà passée
	if time.Now().After(activityStartDate) {
		return errors.New("impossible de s'inscrire à une activité passée")
	}

	// Démarrer une transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Inscrire l'utilisateur
	_, err = tx.Exec(
		"INSERT INTO registrations (user_id, activity_id) VALUES (?, ?)",
		userID, activityID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Valider la transaction
	return tx.Commit()
}

// UnregisterFromActivity désinscrire un utilisateur d'une activité
func UnregisterFromActivity(db *sql.DB, userID, activityID int64) error {
	// Vérifier si l'utilisateur est inscrit
	var isRegistered bool
	err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM registrations WHERE user_id = ? AND activity_id = ?)",
		userID, activityID,
	).Scan(&isRegistered)

	if err != nil {
		return err
	}

	if !isRegistered {
		return errors.New("vous n'êtes pas inscrit à cette activité")
	}

	// Vérifier si l'activité n'est pas déjà passée
	var activityStartDate time.Time

	err = db.QueryRow(
		"SELECT start_date FROM activities WHERE id = ?",
		activityID,
	).Scan(&activityStartDate)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("activité non trouvée")
		}
		return err
	}

	if time.Now().After(activityStartDate) {
		return errors.New("impossible de se désinscrire d'une activité passée")
	}

	// Désinscrire l'utilisateur
	_, err = db.Exec(
		"DELETE FROM registrations WHERE user_id = ? AND activity_id = ?",
		userID, activityID,
	)

	return err
}

// GetUserRegistrations récupère les activités auxquelles un utilisateur est inscrit
func GetUserRegistrations(db *sql.DB, userID int64, includeHistory bool) ([]models.Activity, error) {
	// Construire la requête
	query := `
		SELECT a.id, a.title, a.description, a.image_path, 
		       a.start_date, a.end_date, a.location, 
		       a.max_participants, a.eco_points, a.created_at, a.updated_at,
		       COUNT(r2.id) as current_participants,
		       1 as user_registered
		FROM activities a
		JOIN registrations r ON a.id = r.activity_id AND r.user_id = ?
		LEFT JOIN registrations r2 ON a.id = r2.activity_id
	`

	// Ajouter les filtres
	whereClause := ""
	var args []interface{}
	args = append(args, userID)

	if !includeHistory {
		whereClause = " WHERE a.end_date >= ?"
		args = append(args, time.Now())
	}

	// Grouper et ordonner
	groupAndOrder := `
		GROUP BY a.id
		ORDER BY a.start_date ASC
	`

	// Requête finale
	query = query + whereClause + groupAndOrder

	// Exécuter la requête
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parcourir les résultats
	activities := []models.Activity{}
	for rows.Next() {
		var activity models.Activity
		var startDate, endDate, createdAt, updatedAt time.Time
		var userRegistered bool

		err := rows.Scan(
			&activity.ID, &activity.Title, &activity.Description, &activity.ImagePath,
			&startDate, &endDate, &activity.Location,
			&activity.MaxParticipants, &activity.EcoPoints, &createdAt, &updatedAt,
			&activity.CurrentParticipants, &userRegistered,
		)

		if err != nil {
			return nil, err
		}

		// Convertir les dates
		activity.StartDate = startDate
		activity.EndDate = endDate
		activity.CreatedAt = createdAt
		activity.UpdatedAt = updatedAt
		activity.UserRegistered = userRegistered

		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return activities, nil
}
