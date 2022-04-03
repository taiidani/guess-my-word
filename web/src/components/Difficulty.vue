<template>
  <div>
    <select id="mode" class="form-select" v-on:change="modeChange">
      <option
        v-for="list in lists"
        :key="list.name"
        v-bind:selected="mode == list.name"
      >
        {{ list.name }}
      </option>
    </select>
  </div>
</template>

<script>
export default {
  name: "Difficulty",
  props: ["mode", "lists"],
  methods: {
    modeChange: modeChange,
  },
};

function modeChange(evt) {
  window.location.replace("/?mode=" + evt.target.value);
}
</script>
