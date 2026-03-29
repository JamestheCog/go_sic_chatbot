// A file to store items and functions pertaining to our PWA's service workers.  Note that 
// I'll cache everything as our application is a pretty tiny one all things considered - 
// ITEMS_TO_CACHE is going to be one long array...

const CACHE_NAME = 'sic_chatbot';
const ITEMS_TO_CACHE = [
    // Essential items to cache first:
    '/', 
    "/manifest.json",

    
    // Logos:
    '/static/img/maskable-512.png',
    '/static/img/192.png',
    '/static/img/512.png',
    '/static/img/logo.png',
    
    // CSS files:
    '/static/css/chat/animations.css',
    '/static/css/chat/header.css',
    '/static/css/chat/main.css',
    '/static/css/chat/voice_input.css',
    '/static/css/chat/wrapper_and_sidebar.css',
    '/static/css/home/main.css',
    '/static/css/login/main.css',
    '/static/css/signup/main.css',
    '/static/css/shared/alerts.css',
    
    // JS files:
    '/static/js/chat/additional_inputs.js',
    '/static/js/chat/audio_input.js',
    '/static/js/chat/main.js',
    '/static/js/chat/sidebar.js',
    '/static/js/chat/theme.js',
    '/static/js/home/main.js',
    '/static/js/login/main.js',
    '/static/js/signup/main.js',
    '/static/js/shared/alerts.js',
    '/static/js/shared/theme_toggler.js'
]

// Installation UI:
self.addEventListener('install', (e) => {
    e.waitUntil(
        caches.open(CACHE_NAME).then(cache => {
            return cache.addAll(ITEMS_TO_CACHE);
        })
    );
    self.skipWaiting(); 
});

// Activation event:
self.addEventListener('activate', (e) => {
    console.log('Service worker activated.');
    event.waitUntil(clients.claim()); 
});

// Fetch event for offline UI:
self.addEventListener('fetch', (e) => {
    e.respondWith(
        caches.match(e.request).then(response => {
            return response || fetch(e.request);
        })
    );
});