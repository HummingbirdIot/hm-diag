<template>
  <CellGroup>
    <Cell>
      <div class="big-title">{{ label }}</div>
    </Cell>
    <Cell>
      <Field v-model="data" autosize type="textarea" readonly class="text-con"></Field>
    </Cell>
  </CellGroup>
</template>

<script setup>
import { ref, onMounted, defineProps } from 'vue';
import { Field, Cell, CellGroup, Toast } from 'vant';

const props = defineProps({
  api: String,
  label: String
})
const { api, label } = props
console.log("api....", props.api)

const data = ref("")

function fetchData() {
  Toast.loading({
    message: '加载中...',
    forbidClick: true,
    loadingType: 'spinner',
    duration: 10 * 1000
  });
  fetch(api)
    .then(r => r.json())
    .then(r => {
      Toast.clear()
      data.value = JSON.stringify(r, null, 2)
    })
    .catch(r => {
      Toast.fail("error :" + r)
    })
}

onMounted(() => {
  fetchData()
})

</script>

<style scoped>
.text-con {
  margin-bottom: 50px;
}
</style>

