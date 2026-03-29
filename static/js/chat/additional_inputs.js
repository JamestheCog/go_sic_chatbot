const IMAGE_UPLOAD = document.querySelector('.attach-button');
const IMAGE_PREVIEW = document.querySelector('.image-preview');
const IMAGE_INPUT = document.querySelector('#imageInput');
const MAX_IMAGE_SIZE = 1024**2 * 10;

// -- Begin helper functions here --
let clearImagePreview = () => {
    IMAGE_PREVIEW.innerHTML = '';
    IMAGE_PREVIEW.value = "";
    IMAGE_INPUT.value = "";
    updateButtonState();
}


// Given an image upload, return its base 64 string and its mime type:
let findImgParams = (img) => {
    if (!img) return;
    console.log(img.type)

    let imgReader = new Promise((resolve, reject) => {
        let fileReader = new FileReader();
        fileReader.onload = function() {
            let result = fileReader.result;
            let b64String = result.split(',')[1];
            resolve({b64_string: b64String, img_mime: img.type});
        };
        fileReader.onerror = reject;
        fileReader.readAsDataURL(img);
    })
    return imgReader;
}


// --- Event listeners -- -
//
// Make it so that when a user does upload their image to the application,
// they'll be able to preview it right next to that clippy icon:
IMAGE_INPUT.addEventListener('change', (e) => {
    IMAGE_PREVIEW.innerHTML = '';
    let file = e.target.files[0];
    if (!file) return;
    if (file.size > MAX_IMAGE_SIZE) {
        showError("The image's too big to be sent!");
        return;
    }

    let wrapper = document.createElement('div');
    wrapper.style.position = 'relative';
    wrapper.style.display = 'inline-block';

    let imagePreview = document.createElement('img');
    imagePreview.src = URL.createObjectURL(file);
    imagePreview.style.maxWidth = '50px';
    imagePreview.style.maxHeight = '50px';
    imagePreview.style.borderRadius = '5px';
    imagePreview.onload = function() {URL.revokeObjectURL(imagePreview.src)};

    // Then, also add in the removal button:
    let removeButton = document.createElement('button');
    removeButton.className = 'remove-button';
    removeButton.textContent = 'x';
    removeButton.addEventListener('click', clearImagePreview);

    // And last but not least, add in the image as a thumbnail-like
    // thin in the IMAGE_PREVIEW container:
    wrapper.appendChild(imagePreview);
    wrapper.appendChild(removeButton);
    IMAGE_PREVIEW.appendChild(wrapper);
});

// Add a click event to that clippy icon: when somebody clicks on it,
// click on the hidden image upload button:
IMAGE_UPLOAD.addEventListener('click', () => {IMAGE_INPUT.click()});