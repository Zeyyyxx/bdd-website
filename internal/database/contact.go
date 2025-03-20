package database

import (
	"database/sql"
	"errors"
	"time"

	"bdd-website/internal/models"
)

// CreateContactMessage enregistre un nouveau message de contact
func CreateContactMessage(db *sql.DB, message models.ContactMessageCreate) (int64, error) {
	// Vérifier que les champs obligatoires sont remplis
	if message.Name == "" || message.Email == "" || message.Subject == "" || message.Message == "" {
		return 0, errors.New("tous les champs sont obligatoires")
	}

	// Insérer le message
	result, err := db.Exec(
		"INSERT INTO contact_messages (name, email, subject, message) VALUES (?, ?, ?, ?)",
		message.Name, message.Email, message.Subject, message.Message,
	)

	if err != nil {
		return 0, err
	}

	// Récupérer l'ID généré
	return result.LastInsertId()
}

// GetContactMessages récupère les messages de contact avec pagination
func GetContactMessages(db *sql.DB, page, pageSize int, unreadOnly bool) ([]models.ContactMessage, int, int, error) {
	// Calculer l'offset pour la pagination
	offset := (page - 1) * pageSize

	// Construire la requête
	query := "SELECT id, name, email, subject, message, submitted_at, is_read FROM contact_messages"
	countQuery := "SELECT COUNT(*) FROM contact_messages"
	unreadCountQuery := "SELECT COUNT(*) FROM contact_messages WHERE is_read = 0"

	// Ajouter les filtres
	whereClause := ""
	if unreadOnly {
		whereClause = " WHERE is_read = 0"
	}

	// Ajouter la clause ORDER BY et LIMIT
	query = query + whereClause + " ORDER BY submitted_at DESC LIMIT ? OFFSET ?"

	// Si on filtre, ajouter la clause WHERE au countQuery
	if unreadOnly {
		countQuery = countQuery + whereClause
	}

	// Exécuter la requête de comptage
	var total int
	err := db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, 0, err
	}

	// Compter les messages non lus
	var unreadCount int
	err = db.QueryRow(unreadCountQuery).Scan(&unreadCount)
	if err != nil {
		return nil, 0, 0, err
	}

	// Exécuter la requête principale
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	// Parcourir les résultats
	messages := []models.ContactMessage{}
	for rows.Next() {
		var message models.ContactMessage
		var submittedAt time.Time

		err := rows.Scan(
			&message.ID, &message.Name, &message.Email, &message.Subject,
			&message.Message, &submittedAt, &message.IsRead,
		)

		if err != nil {
			return nil, 0, 0, err
		}

		message.SubmittedAt = submittedAt
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	return messages, total, unreadCount, nil
}

// GetContactMessage récupère un message de contact spécifique
func GetContactMessage(db *sql.DB, messageID int64) (*models.ContactMessage, error) {
	var message models.ContactMessage
	var submittedAt time.Time

	err := db.QueryRow(
		"SELECT id, name, email, subject, message, submitted_at, is_read FROM contact_messages WHERE id = ?",
		messageID,
	).Scan(
		&message.ID, &message.Name, &message.Email, &message.Subject,
		&message.Message, &submittedAt, &message.IsRead,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("message non trouvé")
		}
		return nil, err
	}

	message.SubmittedAt = submittedAt

	return &message, nil
}

// MarkContactMessageAsRead marque un message de contact comme lu
func MarkContactMessageAsRead(db *sql.DB, messageID int64) error {
	// Vérifier si le message existe
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM contact_messages WHERE id = ?)", messageID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("message non trouvé")
	}

	// Marquer comme lu
	_, err = db.Exec(
		"UPDATE contact_messages SET is_read = 1 WHERE id = ?",
		messageID,
	)

	return err
}

// DeleteContactMessage supprime un message de contact
func DeleteContactMessage(db *sql.DB, messageID int64) error {
	// Vérifier si le message existe
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM contact_messages WHERE id = ?)", messageID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("message non trouvé")
	}

	// Supprimer le message
	_, err = db.Exec(
		"DELETE FROM contact_messages WHERE id = ?",
		messageID,
	)

	return err
}
