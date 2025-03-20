// Charger les informations du profil
async function loadUserProfile() {
    try {
        const response = await fetch('/api/users/profile', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });

        if (!response.ok) {
            throw new Error('Impossible de charger le profil');
        }

        const profile = await response.json();
        
        // Remplir les champs du formulaire
        document.getElementById('username').value = profile.username;
        document.getElementById('email').value = profile.email;

        // Mettre à jour le tableau de bord écologique
        document.getElementById('total-points').textContent = profile.total_eco_points;
        document.getElementById('activity-count').textContent = profile.activity_count;
        document.getElementById('badge-count').textContent = profile.badge_count;
    } catch (error) {
        console.error('Erreur de chargement du profil:', error);
        alert('Impossible de charger les informations du profil');
    }
}

// Charger le tableau de bord écologique
async function loadEcoDashboard() {
    try {
        const [summaryResponse, activitiesResponse] = await Promise.all([
            fetch('/api/eco-dashboard/summary', {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            }),
            fetch('/api/users/registrations', {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            })
        ]);

        if (!summaryResponse.ok || !activitiesResponse.ok) {
            throw new Error('Impossible de charger le tableau de bord');
        }

        const summary = await summaryResponse.json();
        const activitiesData = await activitiesResponse.json();

        // Mettre à jour le classement
        if (summary.ranking && summary.total_users) {
            document.getElementById('user-ranking').textContent = 
                `${summary.ranking}/${summary.total_users}`;
        }

        // Afficher les activités
        const activitiesContainer = document.getElementById('user-activities');
        const noActivitiesMessage = document.getElementById('no-activities');

        if (activitiesData.activities.length === 0) {
            noActivitiesMessage.classList.remove('hidden');
        } else {
            noActivitiesMessage.classList.add('hidden');
            
            activitiesData.activities.forEach(activity => {
                const activityCard = document.createElement('div');
                activityCard.className = 'bg-gray-100 rounded-lg p-4 shadow-sm';
                activityCard.innerHTML = `
                    <h3 class="font-bold text-lg mb-2">${activity.title}</h3>
                    <p class="text-gray-600 mb-2">${activity.location}</p>
                    <div class="flex justify-between">
                        <span class="text-sm text-green-600">${new Date(activity.start_date).toLocaleDateString()}</span>
                        <span class="text-sm text-blue-600">${activity.eco_points} points</span>
                    </div>
                `;
                activitiesContainer.appendChild(activityCard);
            });
        }
    } catch (error) {
        console.error('Erreur de chargement du tableau de bord:', error);
    }
}

// Mettre à jour le profil
document.getElementById('profile-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const username = document.getElementById('username').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    const updateData = { username, email };
    if (password) {
        updateData.password = password;
    }

    try {
        const response = await fetch('/api/users/profile', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: JSON.stringify(updateData)
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Erreur de mise à jour du profil');
        }

        alert('Profil mis à jour avec succès');
        
        // Recharger les informations du profil
        await loadUserProfile();
    } catch (error) {
        console.error('Erreur de mise à jour du profil:', error);
        alert(error.message);
    }
});

// Charger les données au chargement de la page
document.addEventListener('DOMContentLoaded', async () => {
    // Vérifier l'authentification
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/login';
        return;
    }

    try {
        await Promise.all([
            loadUserProfile(),
            loadEcoDashboard()
        ]);
    } catch (error) {
        console.error('Erreur de chargement initial:', error);
        window.location.href = '/login';
    }
});