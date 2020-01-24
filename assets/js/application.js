require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");

let state = new Object();
state.before = [];
state.after = [];
state.start = new Date();
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

    $.get("/guess?word=" + word, function (data) {
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
        $("#guess-box").text("🎉 You guessed \"" + state.answer + "\" correctly in " + state.guesses + " tries. Come back tomorrow for another!");
    }
}
