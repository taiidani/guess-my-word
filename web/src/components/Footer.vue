<template>
  <footer class="container text-center text-md-start px-4">
    <div class="row gx-5">
      <div class="col">
        <div class="p-3">
          <h5><i class="bi bi-clock-history"></i> Yesterday's Stats</h5>
          <stats v-bind:day="yesterday"></stats>
        </div>
      </div>
      <div class="col">
        <div class="p-3">
          <h5><i class="bi bi-bar-chart-fill"></i> Today's Stats</h5>
          <stats v-bind:day="today"></stats>
        </div>
      </div>
      <div class="col">
        <div class="p-3">
          <h5><i class="bi bi-speedometer"></i> List</h5>
          <difficulty v-bind:mode="mode" v-bind:lists="lists" />
        </div>
      </div>
      <div class="col">
        <div class="p-3">
          <h5><i class="bi bi-github"></i> Source</h5>
          <p>
            <a href="https://github.com/taiidani/guess-my-word/releases"
              >Changelog</a
            >
          </p>
          <p>
            <a href="https://github.com/taiidani/guess-my-word"
              >View on GitHub</a
            >
          </p>
        </div>
      </div>
    </div>
  </footer>
</template>

<script>
import Stats from "./Stats.vue";
import Difficulty from "./FooterMode.vue";

let yesterday = {
  word: null,
  completions: 0,
  bestRun: 0,
  avgRun: 0,
};

let today = {
  completions: 0,
  bestRun: 0,
  avgRun: 0,
};

let lists = [];

export default {
  components: { Stats, Difficulty },
  name: "Footer",
  props: ["mode"],
  data() {
    window.setInterval(() => {
      refreshStats(this.mode);
    }, 60000);
    refreshStats(this.mode);

    // Populate the lists dropdown
    fetch("/api/lists")
      .then((response) => response.json())
      .then((data) => {
        console.debug(data);

        if (data.error) {
          console.error(data.error);
          return;
        }

        data.forEach((item) => {
          lists.push(item);
        });
      })
      .catch((err) => {
        console.error(err);
      });

    return {
      yesterday: yesterday,
      today: today,
      lists: lists,
    };
  },
};

function refreshStats(mode) {
  const dt = new Date();
  const params = new URLSearchParams({
    date: Math.floor(dt.getTime() / 1000) - 24 * 60 * 60, // Subtract 1 day
    tz: dt.getTimezoneOffset(),
    mode: mode,
  });

  fetch("/api/stats?" + params.toString())
    .then((response) => response.json())
    .then((data) => {
      console.debug(data);

      if (data.error) {
        console.error(data.error);
        return;
      }

      var stats = analyzeStats(data.word.guesses);
      yesterday.word = data.word.value;
      yesterday.completions = stats.completions;
      yesterday.bestRun = stats.bestRun;
      yesterday.avgRun = stats.avgRun;

      stats = analyzeStats(data.today.guesses);
      today.completions = stats.completions;
      today.bestRun = stats.bestRun;
      today.avgRun = stats.avgRun;
    })
    .catch((err) => {
      console.error(err);
    });
}

function analyzeStats(guesses) {
  if (guesses == null || guesses.length == 0) {
    return {
      completions: 0,
      bestRun: 0,
      avgRun: 0,
    };
  }
  var guessCount = 0;
  var bestRun = 999;
  guesses.forEach((item) => {
    guessCount += item.count;
    if (item.count < bestRun) {
      bestRun = item.count;
    }
  });
  var avgRun = guessCount / guesses.length;
  return {
    completions: guesses.length,
    bestRun: bestRun,
    avgRun: avgRun,
  };
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
