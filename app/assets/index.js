document.addEventListener('htmx:afterSwap', function (evt) {
    let fInput = document.getElementById("guess-input");
    if (fInput !== null) {
        fInput.focus();
    }
});


document.addEventListener("htmx:responseError", function (evt) {
    if (evt.detail.elt.id != "guess-form") {
        return;
    }

    // Set the guess input to an invalid state
    var input = document.getElementById("guess-input");
    input.ariaInvalid = true;
    document.getElementById("invalid-helper").innerHTML = evt.detail.xhr.responseText;
});

// Clear any invalid states when typing in the guess input
document.getElementById("guess-input").addEventListener("keydown", function (evt) {
    this.ariaInvalid = null;
    document.getElementById("invalid-helper").innerHTML = "";
})

console.log("Loaded");
