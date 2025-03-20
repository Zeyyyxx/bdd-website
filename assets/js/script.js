// JavaScript pour le site web du Bureau du Développement Durable

document.addEventListener('DOMContentLoaded', function() {
    // Gestion du menu mobile
    const hamburger = document.querySelector('.hamburger');
    const nav = document.querySelector('nav');
    const overlay = document.querySelector('.overlay');
    
    if (hamburger) {
        hamburger.addEventListener('click', function() {
            hamburger.classList.toggle('active');
            nav.classList.toggle('active');
            
            if (overlay) {
                overlay.classList.toggle('active');
            }
            
            // Empêcher le défilement du body quand le menu est ouvert
            document.body.classList.toggle('no-scroll');
        });
    }
    
    if (overlay) {
        overlay.addEventListener('click', function() {
            hamburger.classList.remove('active');
            nav.classList.remove('active');
            overlay.classList.remove('active');
            document.body.classList.remove('no-scroll');
        });
    }
    
    // Fermer le menu si clic sur un lien
    const navLinks = document.querySelectorAll('nav a');
    navLinks.forEach(link => {
        link.addEventListener('click', function() {
            hamburger.classList.remove('active');
            nav.classList.remove('active');
            
            if (overlay) {
                overlay.classList.remove('active');
            }
            
            document.body.classList.remove('no-scroll');
        });
    });
    
    // Formulaires
    const forms = document.querySelectorAll('form');
    forms.forEach(form => {
        form.addEventListener('submit', function(e) {
            // Si le formulaire a une classe 'needs-validation'
            if (form.classList.contains('needs-validation')) {
                if (!form.checkValidity()) {
                    e.preventDefault();
                    e.stopPropagation();
                }
                
                form.classList.add('was-validated');
            }
        });
    });
    
    // Animation pour faire apparaître les éléments au scroll
    const animateElements = document.querySelectorAll('.animate-on-scroll');
    
    if (animateElements.length > 0) {
        // Fonction pour vérifier si un élément est visible dans le viewport
        function isElementInViewport(el) {
            const rect = el.getBoundingClientRect();
            return (
                rect.top <= (window.innerHeight || document.documentElement.clientHeight) * 0.8
            );
        }
        
        // Fonction pour animer les éléments visibles
        function animateOnScroll() {
            animateElements.forEach(el => {
                if (isElementInViewport(el)) {
                    el.classList.add('animated');
                }
            });
        }
        
        // Exécuter une fois au chargement
        animateOnScroll();
        
        // Et à chaque scroll
        window.addEventListener('scroll', animateOnScroll);
    }
    
    // Gestion des alertes fermables
    const alertCloseButtons = document.querySelectorAll('.alert .close');
    alertCloseButtons.forEach(button => {
        button.addEventListener('click', function() {
            const alert = this.closest('.alert');
            alert.classList.add('fade-out');
            
            setTimeout(() => {
                alert.style.display = 'none';
            }, 300);
        });
    });
    
    // Gestion des onglets
    const tabLinks = document.querySelectorAll('.tab-link');
    tabLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            
            // Désactiver tous les onglets
            const tabLinks = document.querySelectorAll('.tab-link');
            tabLinks.forEach(tab => {
                tab.classList.remove('active');
            });
            
            // Cacher tous les contenus
            const tabContents = document.querySelectorAll('.tab-content');
            tabContents.forEach(content => {
                content.classList.remove('active');
            });
            
            // Activer l'onglet cliqué
            this.classList.add('active');
            
            // Afficher le contenu correspondant
            const targetId = this.getAttribute('data-tab');
            document.getElementById(targetId).classList.add('active');
        });
    });
    
    // Initialisation des tooltips
    const tooltips = document.querySelectorAll('[data-tooltip]');
    tooltips.forEach(tooltip => {
        tooltip.addEventListener('mouseenter', function() {
            const text = this.getAttribute('data-tooltip');
            const tooltipEl = document.createElement('div');
            tooltipEl.className = 'tooltip';
            tooltipEl.textContent = text;
            document.body.appendChild(tooltipEl);
            
            const rect = this.getBoundingClientRect();
            tooltipEl.style.left = rect.left + (rect.width / 2) - (tooltipEl.offsetWidth / 2) + 'px';
            tooltipEl.style.top = rect.top - tooltipEl.offsetHeight - 10 + 'px';
            
            tooltipEl.classList.add('visible');
            
            this.addEventListener('mouseleave', function() {
                tooltipEl.remove();
            });
        });
    });
    
    // Validation côté client pour les formulaires d'inscription et de connexion
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', function(e) {
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            
            if (!email || !password) {
                e.preventDefault();
                showAlert('Veuillez remplir tous les champs.', 'danger');
            }
        });
    }
    
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', function(e) {
            const username = document.getElementById('username').value;
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            const passwordConfirm = document.getElementById('passwordConfirm').value;
            
            if (!username || !email || !password || !passwordConfirm) {
                e.preventDefault();
                showAlert('Veuillez remplir tous les champs.', 'danger');
            } else if (password !== passwordConfirm) {
                e.preventDefault();
                showAlert('Les mots de passe ne correspondent pas.', 'danger');
            }
        });
    }
    
    // Fonction pour afficher des alertes
    function showAlert(message, type = 'info') {
        const alertContainer = document.getElementById('alertContainer');
        if (!alertContainer) return;
        
        const alert = document.createElement('div');
        alert.className = `alert alert-${type}`;
        alert.innerHTML = `${message} <button type="button" class="close">&times;</button>`;
        
        alertContainer.appendChild(alert);
        
        // Fermeture automatique après 5 secondes
        setTimeout(() => {
            alert.classList.add('fade-out');
            setTimeout(() => {
                alert.remove();
            }, 300);
        }, 5000);
        
        // Bouton de fermeture
        const closeButton = alert.querySelector('.close');
        closeButton.addEventListener('click', function() {
            alert.classList.add('fade-out');
            setTimeout(() => {
                alert.remove();
            }, 300);
        });
    }
    
    // Gestion des inscriptions aux activités
    const registerButtons = document.querySelectorAll('.register-activity');
    registerButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            e.preventDefault();
            const activityId = this.getAttribute('data-id');
            
            // Simuler une requête AJAX
            // En production, remplacer par un vrai appel API
            setTimeout(() => {
                this.textContent = 'Inscrit';
                this.classList.remove('btn-primary');
                this.classList.add('btn-success');
                this.disabled = true;
                
                showAlert('Vous êtes inscrit à cette activité!', 'success');
            }, 500);
        });
    });
    
    // Compteur pour des statistiques (animation)
    const counters = document.querySelectorAll('.counter');
    if (counters.length > 0) {
        counters.forEach(counter => {
            const target = +counter.getAttribute('data-target');
            const duration = 1500; // ms
            const steps = 50;
            const stepTime = duration / steps;
            const stepValue = target / steps;
            let current = 0;
            
            function updateCounter() {
                current += stepValue;
                if (current < target) {
                    counter.textContent = Math.round(current);
                    setTimeout(updateCounter, stepTime);
                } else {
                    counter.textContent = target;
                }
            }
            
            // Démarrer le compteur quand il est visible
            const observer = new IntersectionObserver((entries) => {
                entries.forEach(entry => {
                    if (entry.isIntersecting) {
                        updateCounter();
                        observer.unobserve(entry.target);
                    }
                });
            }, { threshold: 0.5 });
            
            observer.observe(counter);
        });
    }
});