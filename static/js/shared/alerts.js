
const ALERT_CONTAINER = document.querySelector('.alert-container');
const TIMEOUT_DURATION = 4500;

// Given a message, make it a success message in the ALERT_CONTAINER element:
let showSuccess = (msg, shouldAutoClose = true) => {
    if (!msg || typeof(msg) !== "string" || msg.trim().length === 0) return;

    let alertDiv = document.createElement('div');
    alertDiv.className = 'alert-success';
    alertDiv.innerHTML = `
        <span class = "message"> ${msg} </span>
        <button class = "close-btn" title = "Close"> x </button>
    `;
    alertDiv.querySelector('.close-btn').addEventListener('click', () => {dismissAlert(alertDiv)});
    ALERT_CONTAINER.appendChild(alertDiv);
    if (ALERT_CONTAINER.children.length > 2) {
        dismissAlert(ALERT_CONTAINER.children[0]);
    }
    if (shouldAutoClose) {
        setTimeout(() => {dismissAlert(alertDiv)}, TIMEOUT_DURATION);
    }
}

// Ditto, but for errors:
let showError = (msg, shouldAutoClose = true) => {
    if (!msg || typeof(msg) !== "string" || msg.trim().length === 0) return;

    let errorDiv = document.createElement('div');
    errorDiv.className = 'alert-error';
    errorDiv.innerHTML = `
        <span class = "message"> ${msg} </span>
        <button class = "close-btn" title = "Close"> x </button>
    `
    errorDiv.querySelector('.close-btn').addEventListener('click', () => {dismissAlert(errorDiv)});
    ALERT_CONTAINER.appendChild(errorDiv);
    if (ALERT_CONTAINER.children.length > 2) {
        dismissAlert(ALERT_CONTAINER.children[0]);
    }
    if (shouldAutoClose) {
        setTimeout(() => {dismissAlert(errorDiv)}, TIMEOUT_DURATION);
    }
}

// A helper function for closing them alerts: 
let dismissAlert = (alertDiv) => {
    if (!alertDiv) return;

    alertDiv.style.animation = 'fadeOut 0.4s forwards';
    alertDiv.addEventListener('animationend', () => {alertDiv.remove()}, {once: true})
}