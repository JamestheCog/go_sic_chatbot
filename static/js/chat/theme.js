// A file to contain JavaScript functions that've to do with the toggling on night- / light-mode
// on the user's device.
//
// Note that this function also works with the user's localStorage variable:

// --- Constants --- 
const EVENING_THRESHOLD = 18;
const MORNING_THRESHOLD = 7;
const NIGHT_MODE_PREF_KEY = 'nightModeOn';
const THEME_BUTTON = document.querySelector('.theme-toggle');

// --- Callbacks ---
//
// Sets the theme for the user based on the state of the application's 
// theme - this function will also update localStorage accordingly:
let toggleTheme = () => {
    const isDark = document.body.classList.toggle('dark');
    THEME_BUTTON.textContent = isDark ? '☀️' : '🌙';
    localStorage.setItem(NIGHT_MODE_PREF_KEY, isDark);
}

// --- Attaching event listeners and callbacks ---
document.addEventListener('DOMContentLoaded', () => {
    let nightModePref = localStorage.getItem(NIGHT_MODE_PREF_KEY);
    let shouldBeDark;

    if (nightModePref !== null) {
        shouldBeDark = nightModePref === 'true';
    } else {
        let currentHour = new Date().getHours();
        shouldBeDark = currentHour < MORNING_THRESHOLD || currentHour > EVENING_THRESHOLD;
        localStorage.setItem(NIGHT_MODE_PREF_KEY, shouldBeDark);
    }

    document.body.classList.toggle('dark', shouldBeDark);
    THEME_BUTTON.textContent = shouldBeDark ? '☀️' : '🌙';
});


// Make that toggleTheme function a callback function for the 
// THEME_BUTTON.
THEME_BUTTON.addEventListener('click', toggleTheme);