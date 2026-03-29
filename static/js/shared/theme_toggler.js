// NOTE: 'nightModeOn' comes from /static/js/chat/theme.js - but it's just moved here 
//       for the sake of convenience...
(function() {
    let pref = localStorage.getItem('nightModeOn');
    let hour = new Date().getHours();
    let shouldBeDark = pref === 'true' || (pref === null && (hour < 7 || hour > 19));
    if (shouldBeDark) document.querySelector('body').classList.add('dark');
})();