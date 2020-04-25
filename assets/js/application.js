// Externalize jQuery to CDN
// require("expose-loader?$!expose-loader?jQuery!jquery");

// Default boostrap JS. Disabled as currently not required
// require("bootstrap/dist/js/bootstrap.bundle.js");

let state = new Object();
state.before = []; // after tracks all guesses that the correct word is before
state.after = []; // after tracks all guesses that the correct word is after
state.start = new Date();  // start tracks when the guessing started
state.end = null; // end tracks when the correct guess was made
state.idleTime = 0; // idleTime tracks how long the user was blocked from guessing (network time)
state.guesses = 0; // guesses is the number of guesses that the user has made
state.answer = ""; // answer stores the correct answer once it has been found

$(() => {
    // Validate localstorage and restore if from the same day
    if (typeof (Storage) === "undefined") {
        alert("Your browser does not appear to support the Local Storage API. Please upgrade to a modern browser in order to activate this feature.");
        return;
    } else if (typeof (sessionStorage.state) !== "undefined") {
        loadState()
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

// loadState will restore the browser's state object from the stored sessionStorage state
function loadState() {
    incomingState = JSON.parse(sessionStorage.state);

    // Massage string dates back into Date objects
    incomingState.start = new Date(incomingState.start);
    if (incomingState.end != null) {
        incomingState.end = new Date(incomingState.end);
    }

    // Default idletime
    if (typeof (incomingState.idleTime) == "undefined") {
        incomingState.idleTime = 0
    }

    console.debug(incomingState);

    // Only assign the state if it started on the same day as the current word
    if (incomingState.start.getUTCDate() == state.start.getUTCDate()
        && incomingState.start.getUTCMonth() == state.start.getUTCMonth()) {
        state = incomingState;
    }
}

function guess(word) {
    // Validate that we haven't guessed this before
    if (state.before.indexOf(word) >= 0 || state.after.indexOf(word) >= 0) {
        alert("You've guessed this word before!");
        return;
    }

    // Populate and track the request while disabling submissions
    requestStart = new Date()
    $("form#guesser button").attr("disabled", "disabled")
    params = { "word": word, "start": state.start.getTime() }
    $.get("/guess?" + $.param(params))
        .done(function (data) {
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

            // Update the state
            sessionStorage.state = JSON.stringify(state);
            console.debug(state);
        })
        .always(function () {
            // Track network request time
            state.idleTime = state.idleTime + (new Date() - requestStart)

            // Restore the ability to make submissions
            $("form#guesser button").removeAttr("disabled")
        });
}

function renderGuesses() {
    // Empty and repopulate the after/before lists
    let beforeElem = $("#before");
    beforeElem.empty();

    if (state.before.length > 0) {
        state.before.sort();
        state.before.forEach(function (item) {
            $("#gutter .guess .word").text(item);
            $("#gutter .guess").clone().appendTo(beforeElem);
        });
    } else {
        $("#gutter .guess .word").html("No guesses after the word");
        $("#gutter .guess").clone().addClass("placeholder").appendTo(beforeElem);
    }

    let afterElem = $("#after");
    afterElem.empty();

    if (state.after.length > 0) {
        state.after.sort();
        state.after.forEach(function (item) {
            $("#gutter .guess .word").text(item);
            $("#gutter .guess").clone().appendTo(afterElem);
        });
    } else {
        $("#gutter .guess .word").html("No guesses before the word");
        $("#gutter .guess").clone().addClass("placeholder").appendTo(afterElem);
    }

    // scroll screen to last after, if available
    scrollElem = $("#after li:nth-last-child(2)").get(0)
    if (typeof scrollElem != "undefined") {
        scrollElem.scrollIntoView()
    }

    // Completion text. Congratulations!
    if (state.answer != "") {
        guessSeconds = Math.floor((state.end - state.start - state.idleTime) / 1000)
        guessMinutes = Math.floor(guessSeconds / 60)
        guessSeconds = guessSeconds % 60
        $("#guesser").text("🎉 You guessed \"" + state.answer + "\" correctly with " + state.guesses + " tries in " + guessMinutes + " minutes, " + guessSeconds + " seconds. Come back tomorrow for another!");
    }
}
