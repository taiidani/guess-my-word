require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");

let state = new Object();
state.before = [];
state.after = [];
state.start = new Date();
state.end = null;
state.guesses = 0;
state.answer = "";

$(() => {
    // Validate localstorage and restore if from the same day
    if (typeof (Storage) === "undefined") {
        alert("Your browser does not appear to support the Local Storage API. Please upgrade to a modern browser in order to activate this feature.");
        return;
    } else if (typeof (sessionStorage.state) !== "undefined") {
        incomingState = JSON.parse(sessionStorage.state);
        incomingState.start = new Date(incomingState.start);

        if (incomingState.end != null) {
            incomingState.end = new Date(incomingState.end);
        }

        console.debug(incomingState);

        if (incomingState.start.getUTCDate() == state.start.getUTCDate()
            && incomingState.start.getUTCMonth() == state.start.getUTCMonth()) {
            state = incomingState;
        }
    }

    // Attach the guess event
    $("form#guesser").submit(function (evt) {
        event.preventDefault();

        let word = $("input", evt.target).first().val();
        $("input", evt.target).first().val("");

        guess(word);
    });

    renderGuesses();
});

function guess(word) {
    // Validate that we haven't guessed this before
    if (state.before.indexOf(word) >= 0 || state.after.indexOf(word) >= 0) {
        alert("You've guessed this word before!");
        return;
    }

    params = { "word": word, "start": state.start.getTime() }
    $.get("/guess?" + $.param(params), function (data) {
        console.debug(data);

        if (data.error != "") {
            alert(data.error);
            return;
        }

        state.guesses += 1;
        if (data.after) {
            state.after.push(word);
        } else if (data.before) {
            state.before.push(word);
        } else if (data.correct) {
            state.answer = word;
            state.end = new Date();
        }

        renderGuesses();
        sessionStorage.state = JSON.stringify(state);
        console.debug(state);
    });
}

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

    if (state.answer != "") {
        guessSeconds = Math.floor((state.end - state.start) / 1000)
        guessMinutes = Math.floor(guessSeconds / 60)
        guessSeconds = guessSeconds % 60
        $("#guess-box").text("🎉 You guessed \"" + state.answer + "\" correctly with " + state.guesses + " tries in " + guessMinutes + " minutes, " + guessSeconds + " seconds. Come back tomorrow for another!");
    }
}
