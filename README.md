# Site Web BDD

Backend Go pour un site web axé sur l'écologie, avec inscription des utilisateurs, activités, défis et tableau de bord des éco-points.

## Fonctionnalités clés

- Inscription et authentification des utilisateurs avec JWT
- Opérations CRUD pour les activités et les défis
- Participation des utilisateurs aux activités et aux défis
- Système d'éco-points et récompenses sous forme de badges
- Tableau de bord administrateur pour gérer le contenu et les utilisateurs
- Formulaire de contact pour la communication avec les utilisateurs

## Stack technologique

- Go 1.23.0
- Base de données SQLite3
- Routeur Gorilla Mux
- JWT pour l'authentification

## Installation

1. Clonez le dépôt
2. Installez les dépendances : `go mod download`
3. Compilez le projet : `go build`
4. Lancez le serveur : `./bdd-website`

Le serveur démarrera sur le port 8080 par défaut. Vous pouvez changer le port et d'autres configurations dans `config/config.go`.

## Structure du projet

- `cmd/` : Point d'entrée principal de l'application
- `config/` : Chargement de la configuration
- `internal/` :
  - `database/` : Connexion à la base de données et requêtes
  - `handlers/` : Gestionnaires HTTP
  - `middleware/` : Middleware personnalisé
  - `models/` : Modèles et structures de données
  - `utils/` : Fonctions utilitaires
- `migrations/` : Scripts de migration de la base de données
- `templates/` : Templates HTML
- `assets/` : Ressources statiques (CSS, JS, images)

## Points d'accès API

Reportez-vous au code dans `internal/handlers/` pour la liste complète des points d'accès API et leurs fonctionnalités.

## Licence

Ce projet est open-source sous la [Licence MIT](LICENSE).
