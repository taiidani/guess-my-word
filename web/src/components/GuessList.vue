<template>
  <div>
    <ul v-bind:class="classes">
      <li class="list-group-item guess no-guesses" v-if="!list.length">
        <span>No guesses {{ name }} the word</span>
      </li>
      <li
        class="list-group-item guess"
        v-for="word in formattedList"
        :key="word.found + word.remaining"
      >
        <span class="found">{{ word.found }}</span>
        <span>{{ word.remaining }}</span>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: "GuessList",
  props: ["list", "name", "known"],
  data() {
    return {
      classes: "list-group " + this.name,
    };
  },
  computed: {
    formattedList: function () {
      var newList = [];
      this.list.forEach((item) => {
        if (item === null) {
          return
        }

        var newItem = { found: "", remaining: "" };

        const searchLength = Math.min(item.length, this.known.length);
        for (var i = 0; i <= searchLength; i++) {
          if (item[i] != this.known[i] || i == searchLength) {
            newItem.found = item.slice(0, i);
            break;
          }
        }

        newItem.remaining = item.slice(newItem.found.length);
        newList.push(newItem);
      });

      return newList;
    },
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

span.found {
  text-decoration: underline;
  background-color: #eef;
}
</style>
