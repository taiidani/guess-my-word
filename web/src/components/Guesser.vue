<template>
  <div class="row">
    <div class="col-md-12">
      <guess-list
        name="before"
        v-bind:list="state.before"
        v-bind:known="known"
      />

      <guess-form v-bind:mode="mode" v-bind:state="state" />

      <guess-list name="after" v-bind:list="state.after" v-bind:known="known" />
    </div>
  </div>
</template>

<script>
import GuessForm from "./GuessForm.vue";
import GuessList from "./GuessList.vue";

export default {
  name: "Guesser",
  components: { GuessForm, GuessList },
  props: ["mode"],
  data() {
    return {
      state: loadState(this.mode),
    };
  },
  computed: {
    // The letters currently known
    known: function () {
      var found = 0;
      if (this.state.after.length == 0 || this.state.before.length == 0) {
        return "";
      }

      var above = this.state.before[this.state.before.length - 1];
      var below = this.state.after[0];

      for (var i = 0; i < Math.max(above.length, below.length); i++) {
        if (above[i] != below[i]) {
          return above.slice(0, i);
        }
      }

      return "";
    },
  },
};

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
    console.log("First time page load; initial state");
    return state;
  }

  // Load from storage
  var incomingState = JSON.parse(sessionStorage["state-" + mode]);

  // Massage string dates back into Date objects
  incomingState.start = new Date(incomingState.start);
  if (incomingState.end != null) {
    incomingState.end = new Date(incomingState.end);
  }

  // Only assign the state if it started on the same day as the current word
  // Having this logic client-side allows a user to keep guessing their word even after
  // a new word has been rotated in. As long as they do not refresh their page their
  // state's "Start" property will lock them to the same day.
  if (
    incomingState.start.getDate() == state.start.getDate() &&
    incomingState.start.getMonth() == state.start.getMonth()
  ) {
    state = incomingState;
  }

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
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
