<template>
  <JsonView :api="api" :label="'Device State'" />
</template>

<script setup>
import { ref } from 'vue';
import JsonView from '../components/JsonView.vue'

const api = ref("/api/v1/device/state")

</script>

