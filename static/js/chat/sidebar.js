// A file to store utilities that've to do with the sidebar on the left:

// Constants 
const SIDEBAR = document.querySelector('.sidebar');
const CHAT_LIST = document.querySelector('.conversation-list');
const SIDEBAR_TOGGLE = document.querySelector('.toggle-sidebar-btn');
const NEW_CONVO = document.querySelector('.new-chat-btn');
const CURRENT_ID_KEY = 'current_conversation_id';
const MAIN_CHAT = document.querySelector('.main-chat');
const LOGOUT = document.querySelector('.user-info');
const LOGOUT_DELAY = 8e3;


// -- Functions for sidebar item management -- 

// Given a JSON object from fetchConversationSnippets(), create a 
// <div> element for each we get.
let createSnippet = (snippetInfo) => {
    if (!snippetInfo.conversation_id || !snippetInfo.date_sent) {
        return;
    }

    let snippetDiv = document.createElement('div');
    snippetDiv.className = 'conversation-item';
    snippetDiv.dataset.conversationId = snippetInfo.conversation_id;
    snippetDiv.innerText = snippetInfo.date_sent;
    snippetDiv.addEventListener('click', async () => {
        let currentID = snippetDiv.dataset.conversationId;
        let storedID = localStorage.getItem(CURRENT_ID_KEY);
        if (storedID) {
            let prev = CHAT_LIST.querySelector(`[data-conversation-id="${storedID}"]`);
            if (prev) prev.classList.toggle('active', false);
        }

        let newSnippet = CHAT_LIST.querySelector(`[data-conversation-id="${currentID}"]`);
        newSnippet.classList.toggle('loading', true);
        let messages = await fetchMessagesFor(currentID);
        newSnippet.classList.toggle('loading', false);
        if (!messages || messages?.status?.code !== 200) {
            showError(`Unable to show messages for conversation ID ${currentID}`);
            prev.click();
            prev.classList.toggle('active', true);
            return;
        };
        CHAT_LIST.querySelector(`[data-conversation-id="${currentID}"]`).classList.toggle('active', true)

        localStorage.setItem(CURRENT_ID_KEY, currentID);
        MESSAGE_BOX.innerHTML = '';
        messages.data.forEach(msg => addMessage(msg['message'], msg['role'], msg['img_b64'], msg['img_mimetype']));
    })
    return snippetDiv;
}

// A function that actually loads the snippets of interest into 
// the side menu - meant to be a callback function:
async function loadSnippets() {
    let snippetData = await fetchSnippetInfo();
    
    if (!snippetData || snippetData?.status?.code != 200) {
        let msg = snippetData ? snippetData.status.message : 'The only snippet is the snip snip you had at birth.'
        showError(msg);
        return;
    }
    snippetData['data'].map(createSnippet).forEach(snippetDiv => CHAT_LIST.appendChild(snippetDiv));
}

// Given a snippet, remove it using our new-fangled animations from Gemini:
let removeSnippet = (snippet) => {
    snippet.classList.add('removing');
    setTimeout(() => snippet.remove(), 300);
}

// The callback of interest to attach to LOGOUT:
async function logout(e) {
    e.preventDefault();
    let request = await fetch('/internal/logout', {method: 'POST', credentials: 'include'})
    let response = request.json();

    if (request.ok) {
        showSuccess(`You've been successfully logged out!  Redirecting you shortly...`);
        setTimeout(() => window.location = '/home', 2e3 + (Math.random() * LOGOUT_DELAY));
    } else {
        showError(`Something went wrong while trying to log you out: ${response.status.message}`);
    }
}


// -- Functions for communicating with the backend --- 

// A function for fetching all conversations in our current session 
// given the session ID in the Go backend.  The backend's going to fetch 
// the first 40 characters of all conversations' first bot messages.
async function fetchSnippetInfo() {
    try {
        let response = await fetch('/internal/fetch_snippets',
            {method: 'POST', credentials: 'include'}
        )
        return await response.json();
    } catch(error) {
        console.error(`Serious error happened: ${error}`);
        showErrorAlert(`Can't load your conversations.  Man.  What happened (${error})...`);
        return;
    }
}

// Given a conversation ID, fetch its messages - assume that the session
// ID will persist on the backend:
async function fetchMessagesFor(convoID) {
    if (typeof(convoID) !== String && convoID.trim().length === 0) return;

    let response = await fetch('/internal/fetch_messages', {
        method: 'POST', credentials: 'include',
        body: JSON.stringify({conversation_id: convoID})
    })
    return await response.json();
}

// The event listener of interest for creating a new conversation:
async function newConversation() {
    try { 
        let request = await fetch('/internal/new_chat', 
            {method: 'POST', credentials: 'include'}
        )
        let response = await request.json();

        if (request.ok) {
            let convoData = response['data'];
            localStorage.setItem(CURRENT_ID_KEY, convoData['conversation_id']);
            
            let dataSnippet = createSnippet(convoData);
            CHAT_LIST.prepend(dataSnippet);
        } else {
            showError(`Could not make a new conversation: ${response['status']['message']}`);
        }
    } catch(error) {
        showError(`Something went seriously wrong with the application: ${error}`);
        return;
    }
}

// --- Attach event listeners here: ---
//
// Hide the sidebar when the user clicks on the hider and vice versa:
document.addEventListener('DOMContentLoaded', () => {
    LOGOUT.addEventListener('click', logout);
    SIDEBAR_TOGGLE.addEventListener('click', () => {
        let isCollapsed = SIDEBAR.classList.toggle('collapsed');
        MAIN_CHAT.classList.toggle('sidebar-collapsed', isCollapsed);
    });
    NEW_CONVO.addEventListener('click', () => {
        if (IMAGE_INPUT.files.length || CHATBOX.value != '') {
            shouldErase = confirm(`You've still unsent inputs - are you sure you want to start a new chat?`)
            if (!shouldErase) return;
        }
        
        MESSAGE_BOX.innerHTML = '';
        CHATBOX.value = '';
        IMAGE_INPUT.value = '';
        localStorage.setItem(CURRENT_ID_KEY, '');
        enableInputs(false);
        setTimeout(() => {
            addMessage(WELCOME_MESSAGE, BOT_ROLE);
            enableInputs();
            updateButtonState();
        }, Math.random() * WELCOME_DELAY);
        showSuccess("Chat successfully reset!");
    })
    loadSnippets();
});