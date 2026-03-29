// A file to store functions and constants that've to do with the audio input of our 
// application.

const VOICE_BUTTON = document.querySelector('.voice-button');

let recognition = null;
let isRecording = false;

// If the user's browser supports audio recording, then configure recognition to be
// a SpeechRecognition object.  Otherwise, hide the mic. button.
//
// Note that this functionality only functions for English FOR NOW:
if ('SpeechRecognition' in window || 'webkitSpechRecoginition' in window) {
    const SPEECH_RECOGNITION = window.SpeechRecognition || window.webkitSpechRecoginition;
    recognition = new SPEECH_RECOGNITION();
    recognition.continuous = false;
    recognition.interimResults = false;
    recognition.lang = 'en-US';

    recognition.onresult = (e) => {
        let transcript = '';
        for (let i = e.resultIndex; i < e.results.length ; i++) {
            transcript += e.results[i][0].transcript;
        }
        CHATBOX.value = transcript;
    }

    recognition.onerror = (e) => {
        console.error(`Error while recording speech: ${e}`);
        stopRecording();
        // alert here:
        showErrorAlert(e);
    }

    recognition.onend = () => {
        stopRecording();
    }
} else {
    VOICE_BUTTON.style.display = 'none';
}


// --- Functions for recording ---
//
// Starts the recording:
let startRecording = () => {
    if (!recognition || isRecording) return;
    isRecording = true;
    VOICE_BUTTON.classList.add('recording');
    recognition.start();
};

// Stops the recording:
let stopRecording = () => {
    if (!isRecording) return
    isRecording = false;
    VOICE_BUTTON.classList.remove('recording');
    recognition.stop();
};


// --- Functionality for recording ---
let holdTimeout;
VOICE_BUTTON.addEventListener('mousedown', startRecording);
VOICE_BUTTON.addEventListener('touchstart', (e) => {
    e.preventDefault();
    startRecording();
})
VOICE_BUTTON.addEventListener('mouseup', stopRecording);
VOICE_BUTTON.addEventListener('mouseleave', stopRecording);
VOICE_BUTTON.addEventListener('touchend', (e) => {
    e.preventDefault();
    stopRecording();
})