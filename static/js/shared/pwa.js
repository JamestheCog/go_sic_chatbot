// A file to store JavaScript that has to do with the installation of the PWA on the 
// user's device:

const SW_PATH = '/sw.js'

if ('serviceWorker' in navigator) {
    window.addEventListener('load', () => {
        navigator.serviceWorker.register(SW_PATH)
        .then(() => console.log('Service worker registered'))
        .catch((error) => console.error(`Could not register the worker: ${error}`));
    });
}