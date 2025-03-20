-- Création des tables

-- Table des utilisateurs
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table des activités
CREATE TABLE activities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    image_path TEXT,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    location TEXT NOT NULL,
    max_participants INTEGER DEFAULT 0,
    eco_points INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table des inscriptions aux activités
CREATE TABLE registrations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    activity_id INTEGER NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (activity_id) REFERENCES activities(id) ON DELETE CASCADE,
    UNIQUE(user_id, activity_id)
);

-- Table des défis écologiques
CREATE TABLE eco_challenges (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    points INTEGER NOT NULL,
    duration_days INTEGER NOT NULL,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table des participations aux défis
CREATE TABLE challenge_participants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    challenge_id INTEGER NOT NULL,
    status TEXT NOT NULL, -- 'in_progress', 'completed', 'abandoned'
    joined_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (challenge_id) REFERENCES eco_challenges(id) ON DELETE CASCADE,
    UNIQUE(user_id, challenge_id)
);

-- Table des points écologiques
CREATE TABLE eco_points (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    activity_id INTEGER,
    challenge_id INTEGER,
    points INTEGER NOT NULL,
    description TEXT NOT NULL,
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (activity_id) REFERENCES activities(id) ON DELETE SET NULL,
    FOREIGN KEY (challenge_id) REFERENCES eco_challenges(id) ON DELETE SET NULL
);

-- Table des badges
CREATE TABLE badges (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    image_path TEXT NOT NULL,
    required_points INTEGER NOT NULL,
    category TEXT NOT NULL -- 'participation', 'challenge', 'special'
);

-- Table des badges des utilisateurs
CREATE TABLE user_badges (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    badge_id INTEGER NOT NULL,
    earned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (badge_id) REFERENCES badges(id) ON DELETE CASCADE,
    UNIQUE(user_id, badge_id)
);

-- Table des messages de contact
CREATE TABLE contact_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    subject TEXT NOT NULL,
    message TEXT NOT NULL,
    submitted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN NOT NULL DEFAULT 0
);

-- Insertion des données initiales

-- Création d'un utilisateur administrateur par défaut (mot de passe: admin123)
-- Note: En production, utiliser un mot de passe plus sécurisé et le hacher correctement
INSERT INTO users (email, username, password_hash, is_admin)
VALUES ('admin@example.com', 'Admin', '$2a$10$JPh0PJoNeHwroDfzF6NW6uXZcs.TY4Kz7GQXudCS3KnCYTu/RgzXm', 1);

-- Insertion des badges de base
INSERT INTO badges (name, description, image_path, required_points, category)
VALUES 
    ('Débutant écolo', 'Bienvenue dans la communauté écologique !', '/assets/images/badges/beginner.svg', 0, 'participation'),
    ('Écologiste en herbe', 'Vous avez accumulé 100 points écologiques', '/assets/images/badges/green_starter.svg', 100, 'participation'),
    ('Champion vert', 'Vous avez accumulé 500 points écologiques', '/assets/images/badges/green_champion.svg', 500, 'participation'),
    ('Maître de la durabilité', 'Vous avez accumulé 1000 points écologiques', '/assets/images/badges/sustainability_master.svg', 1000, 'participation'),
    ('Premier défi', 'Vous avez complété votre premier défi', '/assets/images/badges/first_challenge.svg', 50, 'challenge'),
    ('Défieur en série', 'Vous avez complété 5 défis', '/assets/images/badges/serial_challenger.svg', 250, 'challenge'),
    ('Bénévole', 'Vous avez participé à votre première activité', '/assets/images/badges/volunteer.svg', 30, 'participation'),
    ('Ambassadeur BDD', 'Vous avez participé à 10 activités', '/assets/images/badges/ambassador.svg', 300, 'participation');

-- Insertion de quelques défis écologiques
INSERT INTO eco_challenges (title, description, points, duration_days, is_active)
VALUES 
    ('Zéro déchet pendant une semaine', 'Essayez de ne produire aucun déchet non recyclable pendant une semaine entière.', 100, 7, 1),
    ('Transport écologique', 'Utilisez uniquement des transports en commun, vélo ou marche pendant 5 jours consécutifs.', 75, 5, 1),
    ('Réduction d''énergie', 'Réduisez votre consommation d''électricité de 20% pendant 10 jours.', 120, 10, 1),
    ('Alimentation locale', 'Ne consommez que des produits locaux (moins de 100km) pendant 3 jours.', 50, 3, 1);

-- Insertion de quelques activités
INSERT INTO activities (title, description, image_path, start_date, end_date, location, max_participants, eco_points)
VALUES 
    ('Atelier zéro déchet', 'Apprenez à fabriquer vos propres produits ménagers écologiques.', '/assets/images/events/workshop.jpg', 
     datetime('now', '+7 days'), datetime('now', '+7 days', '+3 hours'), 'Salle A103, Paris Ynov Campus', 20, 30),
    
    ('Nettoyage du parc', 'Collecte de déchets dans le parc à proximité du campus.', '/assets/images/events/cleanup.jpg', 
     datetime('now', '+14 days'), datetime('now', '+14 days', '+4 hours'), 'Parc Martin Luther King', 30, 50),
    
    ('Conférence sur l''économie circulaire', 'Venez découvrir comment réduire votre impact environnemental grâce à l''économie circulaire.', '/assets/images/events/conference.jpg', 
     datetime('now', '+21 days'), datetime('now', '+21 days', '+2 hours'), 'Amphithéâtre, Paris Ynov Campus', 100, 20);