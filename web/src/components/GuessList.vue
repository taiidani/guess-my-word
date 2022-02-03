<template>
  <div>
    <ul v-bind:class="classes">
      <li class="list-group-item guess no-guesses" v-if="!list.length">
        <span class="word">No guesses {{ name }} the word</span>
      </li>
      <li class="list-group-item guess" v-for="word in list" :key="word">
        <span class="word">{{ word }}</span>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: "GuessList",
  props: ["list", "name"],
  data() {
    return {
      classes: "list-group " + this.name,
    };
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.list-group {
  font-size: 1em;
  color: #999;
}

.before li:last-child {
  font-size: 2em;
  color: #111;
}

.after li:first-child {
  font-size: 2em;
  color: #111;
}

li.no-guesses {
  font-size: 1em !important;
  color: #999 !important;
}
</style>
