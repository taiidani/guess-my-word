<template>
  <div>
    <div v-if="day.word"><strong>Word:</strong> {{ day.word }}</div>
    <div><strong>Completions:</strong> {{ day.completions }}</div>
    <div><strong>Best Run:</strong> {{ day.bestRun }}</div>
    <div><strong>Average Run:</strong> {{ Math.floor(day.avgRun) }}</div>
  </div>
</template>

<script>
export default {
  name: "Stats",
  props: ["day"],
};
</script>
