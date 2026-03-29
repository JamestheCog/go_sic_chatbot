// A file created to store elements and front-end logic pertaining to the application's
// feature (and home) page:

const NEW_USER_BUTTON = document.querySelector('.main-feature');
const NEW_USER_ROUTE = '/signup';

// --- Add our event-listeners here:

NEW_USER_BUTTON.addEventListener('click', () => {window.location = NEW_USER_ROUTE});
