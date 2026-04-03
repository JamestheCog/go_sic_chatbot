const STATUS = document.querySelector('.status');
const RETRY = document.querySelector('.retry-btn');
const RETRY_TIMEOUT = 1200;

// A callback meant to be used on the RETRY button above:
let retryConnection = () => {
    STATUS.textContent = 'Checking your connection...';
    RETRY.textContent = 'Trying again...';
    RETRY.disabled = true;

    setTimeout(() => {
        if (navigator.online) {
            window.location.reload();
        } else {
            STATUS.textContent = 'Still offline; do check your internet!';
            RETRY.textContent = 'Check connection again';
            RETRY.disabled = false;
        }
    }, RETRY_TIMEOUT)
}

// Attach our event listeners here:
RETRY.addEventListener('click', retryConnection);
window.addEventListener('online', () => {window.location.reload()});