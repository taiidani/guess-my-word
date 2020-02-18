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

    // Attach the complete event
    $("#completion form").submit(function (evt) {
        event.preventDefault();

        let username = $("#username", evt.target).val();
        let suggestion = $("#suggestion", evt.target).val();

        complete(username, suggestion);
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

function complete(username, suggestion) {
    // Populate and track the request while disabling submissions
    $("#completion button").attr("disabled", "disabled")
    params = {
        "word": state.answer,
        "suggestion": suggestion,
        "username": username,
        "start": state.start.getTime(),
        "duration": getTotalDuration(),
        "guesses": state.guesses,
    }
    $.post("/complete", params)
        .done(function (data) {
            $("#completion form").remove()
        })
        .fail(function (data) {
            if (data.responseJSON.error != undefined && data.responseJSON.error != "") {
                alert(data.responseJSON.error);
            }
        })
        .always(function (data) {
            console.debug(data);
            $("#completion button").removeAttr("disabled")
        });
}

function renderGuesses() {
    // Empty and repopulate the after/before lists
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

    // scroll screen to last after, if available
    scrollElem = $("#after li:nth-last-child(2)").get(0)
    if (typeof scrollElem != "undefined") {
        scrollElem.scrollIntoView()
    }

    // Completion text. Congratulations!
    if (state.answer != "") {
        guessSeconds = getTotalDuration()
        guessMinutes = Math.floor(guessSeconds / 60)
        guessSeconds = guessSeconds % 60

        $("#guesser").remove()
        $("#answer").text(state.answer)
        $("#guesses").text(state.guesses)
        $("#minutes").text(guessMinutes)
        $("#seconds").text(guessSeconds);
        $("#completion").show()
        $("#completion input").first().focus()
    }
}

function getTotalDuration() {
    return Math.floor((state.end - state.start - state.idleTime) / 1000)
}