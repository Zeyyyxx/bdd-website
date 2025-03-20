package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initialise la connexion à la base de données et exécute les migrations si nécessaire
func InitDB(dbPath string) (*sql.DB, error) {
	// Vérifier si le fichier de base de données existe
	dbExists := true
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		dbExists = false
	}

	// Créer le répertoire si nécessaire
	dbDir := filepath.Dir(dbPath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("impossible de créer le répertoire de la base de données: %v", err)
		}
	}

	// Ouvrir la connexion à la base de données
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("impossible d'ouvrir la base de données: %v", err)
	}

	// Tester la connexion
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("impossible de se connecter à la base de données: %v", err)
	}

	// Si la base de données n'existait pas, exécuter les migrations
	if !dbExists {
		if err := runMigrations(db); err != nil {
			return nil, fmt.Errorf("erreur lors de l'exécution des migrations: %v", err)
		}
	}

	return db, nil
}

// runMigrations exécute le script SQL de migration initial
func runMigrations(db *sql.DB) error {
	migrationFile := "./migrations/init.sql"

	// Lire le contenu du fichier de migration
	migrationSQL, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("impossible de lire le fichier de migration: %v", err)
	}

	// Exécuter les migrations dans une transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("impossible de démarrer une transaction: %v", err)
	}

	if _, err := tx.Exec(string(migrationSQL)); err != nil {
		tx.Rollback()
		return fmt.Errorf("erreur lors de l'exécution des migrations: %v", err)
	}

	// Commit de la transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("impossible de committer la transaction: %v", err)
	}

	return nil
}
