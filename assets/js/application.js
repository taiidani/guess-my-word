let mode = $("body").data("mode");

// Initialize Vue
let vm = new Vue({
    el: "#guess-area",
    data: {
        state: loadState(),
        word: "",
    },
    computed: {
        // The number of minutes spent guessing
        guessMinutes: function () {
            guessSeconds = Math.floor((this.state.end - this.state.start - this.state.idleTime) / 1000)
            return Math.floor(guessSeconds / 60)
        },
        // The number of seconds spent guessing
        guessSeconds: function () {
            guessSeconds = Math.floor((this.state.end - this.state.start - this.state.idleTime) / 1000)
            return guessSeconds % 60
        },
    },
    methods: {
        guess: guess,
        hint: hint,
        reveal: reveal,
    },
});

// Attach the mode change event
$("#mode").change(function (evt) {
    event.preventDefault();
    window.location.replace("/?mode=" + $("#mode").val());
    return;
});

renderGuesses();

// loadState will restore the browser's state object from the stored sessionStorage state
function loadState() {
    let state = {
        version: 0.9, // Schema version of the state that has been stored
        before: [], // after tracks all guesses that the correct word is before
        after: [], // after tracks all guesses that the correct word is after
        start: new Date(),  // start tracks when the guessing started
        end: null, // end tracks when the correct guess was made
        idleTime: 0, // idleTime tracks how long the user was blocked from guessing (network time)
        guesses: 0, // guesses is the number of guesses that the user has made
        answer: "", // answer stores the correct answer once it has been found
    };

    // Determine if session storage is supported
    if (typeof (Storage) === "undefined") {
        alert("Your browser does not appear to support the Local Storage API. Please upgrade to a modern browser in order to persist guesses across page reloads.");
        return state;
    } else if (typeof (sessionStorage["state-" + mode]) === "undefined") {
        // First-time page load, no state
        return state;
    }

    // Load from storage
    incomingState = JSON.parse(sessionStorage["state-" + mode]);

    // Massage string dates back into Date objects
    incomingState.start = new Date(incomingState.start);
    if (incomingState.end != null) {
        incomingState.end = new Date(incomingState.end);
    }

    // Only assign the state if it started on the same day as the current word
    // Having this logic client-side allows a user to keep guessing their word even after
    // a new word has been rotated in. As long as they do not refresh their page their
    // state's "Start" property will lock them to the same day.
    if (incomingState.start.getDate() == state.start.getDate()
        && incomingState.start.getMonth() == state.start.getMonth()) {
        state = incomingState;
    }

    return state
}

function saveState(state) {
    // Determine if session storage is supported
    if (typeof (Storage) === "undefined") {
        return
    }

    serialized = JSON.stringify(state);
    sessionStorage["state-" + mode] = serialized;
    console.debug(serialized);
}

function guess() {
    let state = this.state
    let word = this.word

    // Validate that we haven't guessed this before
    if (state.before.indexOf(word) >= 0 || state.after.indexOf(word) >= 0) {
        alert("You've guessed this word before!");
        return;
    }

    // Populate and track the request while disabling submissions
    requestStart = new Date()
    $("form#guesser button").attr("disabled", "disabled")
    params = {
        "guesses": state.guesses,
        "word": word,
        "start": Math.floor(state.start.getTime() / 1000),
        "tz": state.start.getTimezoneOffset(),
        "mode": mode,
    }
    $.get("/guess?" + $.param(params))
        .done(function (data) {
            console.debug(data);

            if (data.error) {
                alert(data.error);
                return;
            }

            state.guesses += 1;
            if (data.after) {
                state.after.push(word);
                state.after.sort();
            } else if (data.before) {
                state.before.push(word);
                state.before.sort();
            } else if (data.correct) {
                state.answer = word;
                state.end = new Date();
            }

            renderGuesses();

            // Update the state
            saveState(state);
        })
        .always(function () {
            // Track network request time
            state.idleTime = state.idleTime + (new Date() - requestStart)

            // Restore the ability to make submissions
            $("form#guesser button").removeAttr("disabled")
        });

    this.word = "";
}

function hint() {
    let state = this.state

    // Validate that we have some guesses
    if (state.before.length == 0 || state.after.length == 0) {
        alert("You need to at least guess the before and after first!");
        return
    }

    params = {
        "before": state.before[0],
        "after": state.after[state.after.length - 1],
        "start": Math.floor(state.start.getTime() / 1000),
        "tz": state.start.getTimezoneOffset(),
        "mode": mode,
    }

    $.get("/hint?" + $.param(params))
        .done(function (data) {
            console.debug(data);

            if (data.error) {
                alert(data.error);
                return;
            }

            alert("The word starts with: '" + data.word + "'");
        });
}


function reveal() {
    dt = new Date()

    params = {
        "date": Math.floor(dt.getTime() / 1000) - (24 * 60 * 60), // Subtract 1 day
        "tz": dt.getTimezoneOffset(),
        "mode": mode,
    }

    $.get("/reveal?" + $.param(params))
        .done(function (data) {
            console.debug(data);

            if (data.error) {
                alert(data.error);
                return;
            }

            txt = "Word: " + data.word.Value + "<br/>";
            if (data.word.Guesses != null && data.word.Guesses.length > 0) {
                guessCount = 0;
                bestRun = 999;
                data.word.Guesses.each((item) => {
                    guessCount += item.Count;
                    if (item.Count < bestRun) {
                        bestRun = item.Count;
                    }
                });
                avgCount = guessCount / data.word.Guesses.length;

                txt += "Completions: " + data.word.Guesses.length + "<br/>";
                txt += "Best Run: " + bestRun + "<br/>";
                txt += "Average Run: " + avgCount + "<br/>";
            }

            document.getElementById("reveal").innerHTML = txt;
        });
}

function renderGuesses() {
    // scroll screen to last after, if available
    scrollElem = $("#after li:nth-last-child(2)").get(0)
    if (typeof scrollElem != "undefined") {
        scrollElem.scrollIntoView()
    }
}
