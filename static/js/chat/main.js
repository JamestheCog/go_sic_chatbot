// A file to store functions that've to do with the general functioning of its front-end.

const CHATBOX = document.querySelector('#userInput');
const SUBMIT_BUTTON = document.querySelector('#sendButton');
const RESET_BUTTON = document.querySelector('.clear-button');
const MESSAGE_BOX = document.querySelector('.messages');
const TYPING_INDICATOR = document.querySelector('.typing-indicator');
const USER_ROLE = 'user';
const BOT_ROLE = 'bot';
const WELCOME_DELAY = 2e3;
const WELCOME_MESSAGE = 'Welcome!  Type something to get started...';

(function() {
    let pref = localStorage.getItem(NIGHT_MODE_PREF_KEY);
    let hour = new Date().getHours();
    let shouldBeDark = pref === 'true' || (pref === null && (hour < 7 || hour > 19));
    if (shouldBeDark) document.documentElement.classList.add('dark');
})();

// --- Helper functions for state management ---
let updateButtonState = () => {
    let hasText = CHATBOX.value.trim().length > 0;
    let hasImage = IMAGE_INPUT.files.length > 0;

    IMAGE_UPLOAD.classList.toggle('active', !hasImage);
    if (hasImage) {
        let shouldBeActive = hasText && hasImage;
        IMAGE_UPLOAD.disabled = true;
        IMAGE_UPLOAD.style.cursor = 'not-allowed';
        SUBMIT_BUTTON.classList.toggle('active', shouldBeActive);
        SUBMIT_BUTTON.disabled = !shouldBeActive;
    } else {
        SUBMIT_BUTTON.classList.toggle('active', hasText);
        SUBMIT_BUTTON.disabled = !hasText;
        VOICE_BUTTON.disabled = !hasText;
        IMAGE_UPLOAD.disabled = false;
        IMAGE_UPLOAD.style.cursor = 'pointer';
    }
}

// Allows the application to toggle button statuses - a `true` 
// is an indicator to enable the following buttons:
let enableInputs = (shouldActivate = true) => {
    NEW_CONVO.disabled = shouldActivate
    SUBMIT_BUTTON.classList.toggle('active', shouldActivate);
    SUBMIT_BUTTON.disabled = !shouldActivate;
    SUBMIT_BUTTON.style.cursor = shouldActivate ? 'pointer' : 'disabled';
    VOICE_BUTTON.disabled = !shouldActivate;
    VOICE_BUTTON.classList.toggle('active', shouldActivate);
    VOICE_BUTTON.style.cursor = shouldActivate ? 'pointer' : 'disabled';
    IMAGE_INPUT.value = '';
    IMAGE_UPLOAD.disabled = !shouldActivate;
    IMAGE_UPLOAD.classList.toggle('active', shouldActivate);
    IMAGE_UPLOAD.style.cursor = shouldActivate ? 'pointer' : 'disabled';
}

// --- Sender / session-related functions ---

// Given an message, a base64 image string, and a mimetype, fetch the payload
// from the Go backend:
async function fetchResponse(msg, imgB64 = '', imgMimeType = '') {
    if (!msg || !msg.trim() || typeof(msg) !== 'string') return;
    let conversationID = localStorage.getItem(CURRENT_ID_KEY);
    if (!conversationID) {
        showError(`There's no conversation ID in the backend - could you refresh the application?`);
        return;
    }
    
    try {
        enableInputs(false);
        let response = await fetch('/internal/chat', {
            method: 'POST', credentials: 'include',
            headers: {'Content-Type': 'application/json'},
            body : JSON.stringify({
                session: {
                    conversation_id: conversationID
                },
                message: {
                    msg: msg, role: USER_ROLE,
                    img_b64: imgB64, img_mime: imgMimeType
                }
            })
        });
        return await response.json();
    } catch(error) {
        showError(`Couldn't process the sent message: ${error}`)
    } finally {
        enableInputs();
        updateButtonState();
    }
}

// A helper function for deleting the chat in question.  This function will return immediately 
// it finds that the user has no current ID (since that would imply that there's no conversation
// to currently delete)
async function deleteConversation() {
    let currentConversation = localStorage.getItem(CURRENT_ID_KEY);
    if (!currentConversation) return;
    let resetChoice = confirm('Are you sure you wish to reset the current conversation?');

    if (resetChoice) {
        let request = await fetch('/internal/delete_conversation', {
            method: 'POST', credentials: 'include', 
            body: JSON.stringify({conversation_id: currentConversation})
        });

        let response = await request.json();
        if (request.ok) {
            MESSAGE_BOX.innerHTML = '';
            if (SIDEBAR) {
                let toRemove = SIDEBAR.querySelector(`[data-conversation-id="${currentConversation}"]`);
                if (toRemove) {
                    removeSnippet(toRemove);
                } else {
                    showError(`Could not find snippet with ID ${currentConversation}.`)
                }
                localStorage.setItem(CURRENT_ID_KEY, '');
                showSuccess('Deleted conversation!');
            }
        } else {
            showError(`Failed to restart the conversation session: ${response['status']['message']}`);
        }
        enableInputs();
    } else {
        return;
    }
}

// A callback meant to be used on the SUBMIT_BUTTON element:
async function chat() {
    let msg = CHATBOX.value.trim();
    let img = IMAGE_INPUT.files[0];
    let b64Img, imgMime;
    if ((!msg || !msg && !img) || SUBMIT_BUTTON.disabled) return;
    if (!localStorage.getItem(CURRENT_ID_KEY)) await newConversation();

    // If we've an image being sent, then get its b64 string first before 
    // we go on:
    if (img) {
        let imgData = await findImgParams(img);
        b64Img = imgData['b64_string'];
        imgMime = imgData['img_mime'];
    }

    // Then, wrap the entire thing in a try-except-finally clause for 
    // the message sending bit - we first want to:
    //
    // 1) Update the application's state.
    // 2) Fetch the response from the Go backend.
    // 3) Update the app.'s state again with the backend's payload.
    //
    // Note that we're not enableInputs(false) here as this call 
    // has already been made while fetching the payload in fetchResponse().
    try {
        CHATBOX.value = '';
        clearImagePreview();
        updateButtonState();
        addMessage(msg, USER_ROLE, img, imgMime);
        showTyping();

        let response = await fetchResponse(msg, b64Img, imgMime);
        if (response['status']['code'] === 200) {
            showTyping(false);
            addMessage(response['data']['message'], BOT_ROLE);
        } else {
            let userMessages = MESSAGE_BOX.querySelectorAll('.message.user');
            showError(`Message processing failed: ${response['status']['message']}`);
            userMessages[userMessages.length - 1].remove();
            if (userMessages.length === 1) {
                let toRemove = localStorage.getItem(CURRENT_ID_KEY);
                removeSnippet(CHAT_LIST.querySelector(`[data-conversation-id="${toRemove}"]`))
                localStorage.setItem(CURRENT_ID_KEY, '');
            }
            CHATBOX.value = msg;
            enableInputs();
        }
    } catch(error) {
        showError(`The application's gone bonkers, yo (error: ${error})!`);
    } finally {
        if (TYPING_INDICATOR.style.display == 'block') showTyping(false);
    }
}

// --- END ---

// --- Functions for the application's front-end --- 
//
// Given a message (in a string), a message role, and an image if there is one,
// Add it to the chat interface.
let addMessage = (msg, role, file, mimetype) => {
    if (!msg || msg.length === 0) return;

    let msgDiv = document.createElement('div');
    msgDiv.classList.add('message', role);

    if (msg) {
        let textChild = document.createElement('div');
        textChild.textContent = msg;
        msgDiv.appendChild(textChild);
    }
    if (file && mimetype) {
        let fileChild = document.createElement('img');
        if (typeof(file) === 'string' && file.length > 0) {
            fileChild.src = `data:${mimetype};base64,${file}`;
            console.log(`data:${mimetype};base64,${file}`)
        } else {
            fileChild.src = URL.createObjectURL(file);
        }
        fileChild.onload = function() {URL.revokeObjectURL(fileChild.src)};
        fileChild.className = 'chat-image';
        msgDiv.appendChild(fileChild);
    }
    MESSAGE_BOX.appendChild(msgDiv);
    MESSAGE_BOX.scrollTop = MESSAGE_BOX.scrollHeight;
}

// Shows or hides the typing indicator:
let showTyping = (shouldShow = true) => {
    if (shouldShow) {
        TYPING_INDICATOR.style.display = 'block';
        MESSAGE_BOX.appendChild(TYPING_INDICATOR);
        MESSAGE_BOX.scrollTop = MESSAGE_BOX.scrollHeight;
    } else {
        TYPING_INDICATOR.style.display = 'none';
    }
}


// --- Attaching event listeners here ---
//
// Attach our event listeners here when the page bloody loads:
document.addEventListener('DOMContentLoaded', () => {
    try {
        localStorage.setItem(CURRENT_ID_KEY, '');
        SUBMIT_BUTTON.addEventListener('click', chat);
        CHATBOX.addEventListener('input', updateButtonState);
        CHATBOX.addEventListener('keypress', (e) => {
            if (e.key == 'Enter' && !e.shiftKey) {
                e.preventDefault();
                SUBMIT_BUTTON.click();
            }
        })
        IMAGE_INPUT.addEventListener('input', updateButtonState);
        if (!MESSAGE_BOX.value) {
            enableInputs(false);
            setTimeout(() => {
                addMessage(WELCOME_MESSAGE, BOT_ROLE);
                enableInputs();
                updateButtonState();
            }, Math.random() * WELCOME_DELAY);
        }
        RESET_BUTTON.addEventListener('click', deleteConversation);
    } catch(error) {
        showError(`The application was unable to initialize itself: ${error}`);
    }
})