<template>
  <div class="row">
    <div class="col-md-12">
      <guess-list name="before" v-bind:list="state.before" v-bind:known="known" />

      <guess-form v-bind:mode="mode" v-bind:state="state" @guess="guess" />

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
  props: ["mode", "state"],
  data() {
    return {};
  },
  methods: {
    guess: guess
  },
  computed: {
    // The letters currently known
    known: function () {
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

function guess(data, word) {
  if (word === null) {
    console.error("guessed word is: null");
    return;
  } else if (word === undefined) {
    console.error("guessed word is: undefined");
    return;
  }

  this.state.guesses += 1;
  if (data.after) {
    this.state.before.push(word);
    this.state.before.sort();
  } else if (data.before) {
    this.state.after.push(word);
    this.state.after.sort();
  } else if (data.correct) {
    this.state.answer = word;
    this.state.end = new Date();
  }

  scrollView();
  this.state.save(this.state);
}

function scrollView() {
  // scroll screen to last after, if available
  const matches = document.querySelectorAll(".before li:nth-last-child(2)");
  if (matches.length > 0) {
    matches[0].scrollIntoView();
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
