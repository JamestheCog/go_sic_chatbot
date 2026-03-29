// A file created to store elements and front-end logic pertaining to the application's
// feature (and home) page:

const FORM = document.querySelector('#newUserForm');
const LOGIN_REDIRECT = 2500;
const CHAT_HANDLER = '/chat';

(function() {
    let pref = localStorage.getItem('nightModeOn');
    let hour = new Date().getHours();
    let shouldBeDark = pref === 'true' || (pref === null && (hour < 7 || hour > 19));
    if (shouldBeDark) document.querySelector('body').classList.add('dark');
})();

// --- Callbacks and other utilities --- 
//
// The callback of interest for when we do want to send form data to the backend:
async function sendNewUser(e) {
    e.preventDefault();
    let newUser = new FormData(FORM);

    // Do some validation checks before proceeding:
    if (!newUser.get('username') || !newUser.get('password') || !newUser.get('confirmPassword')) {
        showError("Sorry bub - please fill everything out.");
        return;
    }
    if (newUser.get('password') !== newUser.get('confirmPassword')) {
        showError('Them passwords no be matching.');
        return;
    }

    try {
        let request = await fetch('/internal/new_user', {
            method: 'POST', credentials: 'include',
            body: newUser
        });
        let response = await request.json();

        if (request.ok) {
            showSuccess(response['data']['message'], false);
            setTimeout(() => {window.location = CHAT_HANDLER}, LOGIN_REDIRECT);
        } else {
            showError(`${response['status']['message']}`);
        }
    } catch(error) {
        showError(`What the hell is happening, man (error ${error})`);
    }
}

// Attach the event listener here:
document.addEventListener('DOMContentLoaded', () => {
    FORM.addEventListener('submit', sendNewUser);
})