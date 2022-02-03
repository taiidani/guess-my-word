<template>
  <div>
    <section>
      <form
        id="guesser"
        v-on:submit.prevent="guess"
        v-if="!state.answer"
        v-bind:disabled="guessingInProgress"
      >
        <img
          id="between-logo"
          src="../assets/between.svg"
          alt="Page Break by Arthur Shlain from the Noun Project"
        />

        <label>
          <input
            name="word"
            placeholder="Enter a word here"
            autocomplete="off"
            autocorrect="off"
            autocapitalize="off"
            spellcheck="off"
            v-model="word"
          />
        </label>

        <button class="btn btn-primary" type="submit">Guess</button>
      </form>

      <p v-else>
        🎉 You guessed "{{ state.answer }}" correctly with
        {{ state.guesses }} tries in
        {{ guessMinutes }}
        minutes, {{ guessSeconds % 60 }} seconds. Come back tomorrow for
        another!
      </p>

      <form
        v-on:submit.prevent="hint"
        v-if="
          !state.answer &&
          state.before.length > 0 &&
          state.after.length > 0 &&
          state.guesses > 15
        "
      >
        <button class="btn btn-link btn-sm" type="submit">Need a hint?</button>
      </form>
    </section>
  </div>
</template>

<script>
export default {
  name: "GuessForm",
  props: ["state", "mode"],
  methods: {
    guess: guess,
    hint: hint,
  },
  data() {
    return {
      guessingInProgress: false,
    };
  },
  computed: {
    // The number of minutes spent guessing
    guessMinutes: function () {
      const guessSeconds = Math.floor(
        (this.state.end - this.state.start - this.state.idleTime) / 1000
      );
      return Math.floor(guessSeconds / 60);
    },
    // The number of seconds spent guessing
    guessSeconds: function () {
      const guessSeconds = Math.floor(
        (this.state.end - this.state.start - this.state.idleTime) / 1000
      );
      return guessSeconds % 60;
    },
  },
};

function guess() {
  let state = this.state;
  let word = this.word;

  // Validate that we haven't guessed this before
  if (state.before.indexOf(word) >= 0 || state.after.indexOf(word) >= 0) {
    alert("You've guessed this word before!");
    return;
  }

  // Populate and track the request while disabling submissions
  const requestStart = new Date();
  this.guessingInProgress = true;
  const params = new URLSearchParams({
    guesses: state.guesses,
    word: word,
    start: Math.floor(state.start.getTime() / 1000),
    tz: state.start.getTimezoneOffset(),
    mode: this.mode,
  });

  fetch("/api/guess?" + params.toString())
    .then((response) => response.json())
    .then((data) => {
      console.debug(data);

      if (data.error) {
        alert(data.error);
        return;
      }

      state.guesses += 1;
      if (data.after) {
        state.before.push(word);
        state.before.sort();
      } else if (data.before) {
        state.after.push(word);
        state.after.sort();
      } else if (data.correct) {
        state.answer = word;
        state.end = new Date();
      }

      scrollView();
      state.save(state);
    })
    .catch((err) => {
      console.error(err);
    })
    .finally((info) => {
      // Track network request time
      state.idleTime = state.idleTime + (new Date() - requestStart);

      // Restore the ability to make submissions
      this.guessingInProgress = false;
    });

  this.word = "";
}

function hint() {
  let state = this.state;

  // Validate that we have some guesses
  if (state.before.length == 0 || state.after.length == 0) {
    alert("You need to at least guess the before and after first!");
    return;
  }

  const params = new URLSearchParams({
    before: state.before[0],
    after: state.after[state.after.length - 1],
    start: Math.floor(state.start.getTime() / 1000),
    tz: state.start.getTimezoneOffset(),
    mode: mode,
  });

  fetch("/api/hint?" + params.toString())
    .then((response) => response.json())
    .then((data) => {
      console.debug(data);

      if (data.error) {
        alert(data.error);
        return;
      }

      alert("The word starts with: '" + data.word + "'");
    })
    .catch((err) => {
      console.error(err);
    });
}

function scrollView() {
  // scroll screen to last after, if available
  const matches = document.querySelectorAll(".before li:nth-last-child(2)");
  if (matches.length > 0) {
    matches[0].scrollIntoView();
  }
}
</script>

<style scoped>
section {
  padding: 1em 0;
}

#between-logo {
  height: 3em;
}
</style>
