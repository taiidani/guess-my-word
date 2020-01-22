require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");

state = new Object();
state.before = [];
state.after = [];
state.start = Date.now();
state.guesses = 0;

$(() => {
    // Validate localstorage
    if (typeof (Storage) === "undefined") {
        alert("Your browser does not appear to support the Local Storage API. Please upgrade to a modern browser in order to activate this feature.");
        return;
    } else if (typeof (sessionStorage.state) !== "undefined") {
        state = JSON.parse(sessionStorage.state);
        renderGuesses();
    }

    $("form#guesser").submit(function (evt) {
        event.preventDefault();

        let word = $("input", evt.target).first().val();
        $("input", evt.target).first().val("");

        // Validate that we haven't guessed this before
        if (state.before.indexOf(word) >= 0 || state.after.indexOf(word) >= 0) {
            alert("You've guessed this word before!");
            return;
        }

        $.get("/guess?word=" + word, function (data) {
            console.debug(data);

            if (data.error != "") {
                alert(data.error);
                return;
            }

            state.guesses += 1;
            if (data.after) {
                console.debug("Guess is after word");
                state.after.push(word);
                renderGuesses();
            } else if (data.before) {
                console.debug("Guess is before word");
                state.before.push(word);
                renderGuesses();
            } else if (data.correct) {
                alert("You got it in " + state.guesses + " guesses!");
            }

            sessionStorage.state = JSON.stringify(state);
            console.debug(state);
        });
    });
});

function renderGuesses() {
    console.debug("Rendering...");
    let beforeElem = $("#before");
    beforeElem.empty();
    let afterElem = $("#after");
    afterElem.empty();

    state.before.sort();
    state.before.forEach(function (item) {
        $("#gutter .guess .word").text(item);
        $("#gutter .guess").clone().appendTo(beforeElem);
    });

    state.after.sort();
    state.after.forEach(function (item) {
        $("#gutter .guess .word").text(item);
        $("#gutter .guess").clone().appendTo(afterElem);
    });
}
