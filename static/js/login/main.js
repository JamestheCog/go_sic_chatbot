// A file to contain JavaScript that has to do with the login page of the 
// application:

const LOGIN_FORM = document.querySelector('#loginForm');
const REDIRECT_TIMEOUT = 3500;

// --- Callbacks and other helper functions...
// Callback function for the login functionality:
async function handleLogin(e) {
    e.preventDefault();
    let loginAttempt = new FormData(LOGIN_FORM);
    if (!loginAttempt.get('username') && !loginAttempt.get('password')) {
        showError("You need to enter a username and a password to proceed.")
        return;
    }
    if (!loginAttempt.get('username')) {
        showError('You did not enter a username.');
        return;
    }
    if (!loginAttempt.get('password')) {
        showError('You did not enter a password.');
        return;
    }

    try {
        let request = await fetch('/internal/login', {
            method: 'POST', credentials: 'include',
            body: loginAttempt
        });
        let response = await request.json();

        if (request.ok) {
            showSuccess(response['data']['message']);
            setTimeout(() => {window.location = '/chat'}, REDIRECT_TIMEOUT);
        } else {
            showError(`${response['status']['message']}`);
        }
    } catch(error) {
        showError(`Something is seriously wrong with the application: ${error}`);
        return;
    }
}

// Callback for if the user got re-directed:
document.addEventListener('DOMContentLoaded', () => {
    let urlParams = new URLSearchParams(window.location.search);
    let msgType = urlParams.get('type');
    let msg = urlParams.get('msg');

    if (msgType && msg) showError(msg);
})

// Attach the event listener to the form:
LOGIN_FORM.addEventListener('submit', handleLogin);