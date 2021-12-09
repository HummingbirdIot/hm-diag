<template>
  <JsonView :api="api" :label="'Miner state'" />
</template>

<script setup>
import { ref } from 'vue';
import JsonView from '../components/JsonView.vue'

const api = ref("/api/v1/miner/state?cache=true")

</script>

