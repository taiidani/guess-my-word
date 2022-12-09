import { reactive } from 'vue'

export default reactive({
    loadState: loadState,
    saveState: saveState,
})


// loadState will restore the browser's state object from the stored sessionStorage state
function loadState(mode) {
    let state = {
        version: 0.9, // Schema version of the state that has been stored
        before: [], // after tracks all guesses that the correct word is before
        after: [], // after tracks all guesses that the correct word is after
        start: new Date(), // start tracks when the guessing started
        end: null, // end tracks when the correct guess was made
        idleTime: 0, // idleTime tracks how long the user was blocked from guessing (network time)
        guesses: 0, // guesses is the number of guesses that the user has made
        answer: "", // answer stores the correct answer once it has been found
        save: (state) => {
            saveState(mode, state);
        },
    };

    // Determine if session storage is supported
    if (typeof Storage === "undefined") {
        alert(
            "Your browser does not appear to support the Local Storage API. Please upgrade to a modern browser in order to persist guesses across page reloads."
        );
        return state;
    } else if (typeof sessionStorage["state-" + mode] === "undefined") {
        // First-time page load, no state
        console.log("first time page load for " + mode + "; initial state");
        return state;
    }

    // Load from storage
    var incomingState = JSON.parse(sessionStorage["state-" + mode]);

    // Massage string dates back into Date objects
    incomingState.start = new Date(incomingState.start);
    if (incomingState.end != null) {
        incomingState.end = new Date(incomingState.end);
    }

    state = incomingState;
    state.save = (state) => {
        saveState(mode, state);
    };

    return state;
}

function saveState(mode, state) {
    // Determine if session storage is supported
    if (typeof Storage === "undefined") {
        return;
    }

    const serialized = JSON.stringify(state);
    sessionStorage["state-" + mode] = serialized;
    console.debug(serialized);
}
