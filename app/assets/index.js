document.addEventListener('htmx:afterSwap', function (evt) {
    let fInput = document.getElementById("guess-input");
    if (fInput !== null) {
        fInput.focus();
    }
});

document.getElementById("mode").addEventListener("change", function (evt) {
    let newMode = evt.target.value;
    window.location.href = "/mode/" + encodeURIComponent(newMode);
});

console.log("Loaded");
