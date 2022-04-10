<template>
  <div>
    <select id="mode" class="form-select" v-on:change="modeChange">
      <option
        v-for="list in lists"
        :key="list"
        v-bind:selected="mode == list.toLowerCase()"
      >{{ list }}</option>
    </select>
  </div>
</template>

<script>
export default {
  name: "FooterMode",
  emits: ['modeChange'],
  props: {
    mode: String,
  },
  methods: {
    modeChange: function (evt) { this.$emit("modeChange", evt.target.value.toLowerCase()); },
  },
  data() {
    return {
      lists: []
    }
  },
  mounted() {
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
          this.lists.push(item.name);
        });
      })
      .catch((err) => {
        console.error(err);
      });
  }
};
</script>
