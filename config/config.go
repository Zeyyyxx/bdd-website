package config

import (
	"os"
	"strconv"
)

// Config représente la configuration de l'application
type Config struct {
	// Serveur
	ServerPort int

	// Base de données
	DatabasePath string

	// JWT
	JWTSecret          string
	JWTExpirationHours int
}

// LoadConfig charge la configuration depuis les variables d'environnement ou utilise des valeurs par défaut
func LoadConfig() *Config {
	config := &Config{
		ServerPort:         8080,
		DatabasePath:       "./bdd.db",
		JWTSecret:          "BDDSecretKey", // À remplacer par une clé sécurisée en production
		JWTExpirationHours: 24,
	}

	// Chargement des variables d'environnement si définies
	if port, exists := os.LookupEnv("SERVER_PORT"); exists {
		if p, err := strconv.Atoi(port); err == nil {
			config.ServerPort = p
		}
	}

	if dbPath, exists := os.LookupEnv("DATABASE_PATH"); exists {
		config.DatabasePath = dbPath
	}

	if jwtSecret, exists := os.LookupEnv("JWT_SECRET"); exists {
		config.JWTSecret = jwtSecret
	}

	if jwtExp, exists := os.LookupEnv("JWT_EXPIRATION_HOURS"); exists {
		if exp, err := strconv.Atoi(jwtExp); err == nil {
			config.JWTExpirationHours = exp
		}
	}

	return config
}
