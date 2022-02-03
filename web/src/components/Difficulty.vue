<template>
  <div>
    <select id="mode" class="form-select" v-on:change="modeChange">
      <option value="default" v-bind:selected="mode == 'default'">
        Default
      </option>
      <option value="hard" v-bind:selected="mode == 'hard'">Hard</option>
    </select>
  </div>
</template>

<script>
export default {
  name: "Difficulty",
  props: ["mode"],
  methods: {
    modeChange: modeChange,
  },
};

function modeChange(evt) {
  window.location.replace("/?mode=" + evt.target.value);
}
</script>
