<template>
  <div>
    <template v:if="day.word == 'huh">
      <div><strong>Word:</strong> {{ day.word }}</div>
    </template>
    <div><strong>Completions:</strong> {{ day.completions }}</div>
    <div><strong>Best Run:</strong> {{ day.bestRun }}</div>
    <div><strong>Average Run:</strong> {{ day.avgRun }}</div>
  </div>
</template>

<script>
export default {
  name: "Stats",
  props: ["day"],
};
</script>
