document.addEventListener("htmx:afterSwap", function (evt) {
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
  document.getElementById("invalid-helper").innerHTML =
    evt.detail.xhr.responseText;
  document.getElementById("invalid-helper").classList.remove("hide");
});

// Clear any invalid states when typing in the guess input
document.getElementById("guesser").addEventListener("keydown", function (evt) {
  document.getElementById("guess-input").ariaInvalid = null;
  document.getElementById("invalid-helper").innerHTML = "";
  document.getElementById("invalid-helper").classList.add("hide");
});

console.log("Loaded");
