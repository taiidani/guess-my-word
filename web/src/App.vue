<template>
  <div id="app" class="container">
    <header>
      <h1>
        <a href="/"
          ><img
            src="./assets/logo.svg"
            alt="Question by Gregor Cresnar from the Noun Project, https://thenounproject.com/search/?q=question&i=540041"
        /></a>
        Guess My Word
      </h1>
    </header>

    <div class="row">
      <div class="col-md-12">
        <p>
          This application hosts a daily word game. Guessing using the form
          below will place your guess above or below the box based on its place
          in the alphabet relative to today's word. Every day a new word will
          become guessable. Good luck!
        </p>
      </div>
    </div>

    <hr />

    <Guesser v-bind:mode="mode" />

    <hr />

    <Footer v-bind:mode="mode" />
  </div>
</template>

<script>
import Footer from "./components/Footer.vue";
import Guesser from "./components/Guesser.vue";

export default {
  components: { Footer, Guesser },
  name: "App",
  data() {
    let mode = getMode();
    // Determine the color for the current list
    let request = new URLSearchParams({ name: mode });
    fetch("/api/list?" + request.toString())
      .then((response) => response.json())
      .then((data) => {
        console.debug(data);

        if (data.error) {
          console.error(data.error);
          return;
        }

        document.getElementsByTagName("body")[0].style.backgroundColor =
          "#" + data.color;
      })
      .catch((err) => {
        console.error(err);
      });

    return {
      mode: mode,
    };
  },
};

function getMode() {
  // Extract the mode from the query parameters
  const params = new URLSearchParams(window.location.search);

  if (params.get("mode") == null) {
    params.set("mode", "default");
  }

  return params.get("mode").toLowerCase();
}
</script>

<style>
body {
  background-color: #224;
  color: #eee;
}

body.mode-hard {
  background-color: #422;
  color: #eee;
}

a {
  color: #eee;
  text-decoration: underline;
}

header {
  background-color: #5295de;
  padding: 10px 20px;
  margin-bottom: 1em;
  border-bottom-right-radius: 3px;
  border-bottom-left-radius: 3px;
}

header img {
  display: inline-block;
  max-height: 2em;
}

header h1 {
  display: inline;
  font-size: 30px;
}

hr {
  border-color: #eee;
}
</style>
