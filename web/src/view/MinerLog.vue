<template>
  <CellGroup>
    <Cell title="Recent Time" class="time-sel">
      <template v-for="item in timeArr">
        <Tag
          type="primary"
          plain
          round
          :class="{ 'tag-sel': item.sel }"
          class="tag-tip"
          @click="selTime($event, item.value)"
        >{{ item.label }}</Tag>&nbsp;
      </template>
    </Cell>
    <!-- <Cell title="From Time">
      <input v-model="fromTime" type="datetime-local" />
    </Cell>
    <Cell title="To Time">
      <input v-model="toTime" type="datetime-local" />
    </Cell>-->
    <Cell title>
      <template #title>
        Filter Log:&nbsp;
        <Tag type="primary" plain round @click="filterTxt = 'JSON up'" class="tag-tip">Up link</Tag>&nbsp;
        <Tag type="primary" plain round @click="filterTxt = 'JSON down'" class="tag-tip">Down link</Tag>
      </template>
      <Field v-model="filterTxt" placeholder="input filter text or select filter on the left" />
    </Cell>
    <Cell title>
      <Button type="primary" size="small" plain @click="fullScreen">Full Screen Log</Button> &nbsp;
      <Button type="primary" size="small" @click="query">Query</Button>
    </Cell>
  </CellGroup>

  <!-- <textarea v-model="log" style="width:100%"></textarea> -->
  <!-- <Divider /> -->
  <!-- <Field v-model="log" type="textarea" autosize disabled class="log"></Field> -->
  <pre id="log-con" class="log">{{ log }}</pre>
</template>

<script setup>
import { ref, reactive } from "vue"
import { CellGroup, Cell, Button, Field, Dialog, Divider, Tag, Toast } from "vant"

const timeArr = reactive([
  { label: '10 Minute', value: 10, sel: true },
  { label: '30 Minute', value: 30 },
  { label: '1 hour', value: 60 },
  { label: '3 hour', value: 60 * 3 },
  { label: '6 hour', value: 60 * 6 },
])
const minutesAgo = ref(10)

// const now = new Date()
// const from = new Date(now.getTime() - 10 * 60 * 1000).Format("yyyy-MM-ddTHH:mm:ss")
// const to = new Date().Format("yyyy-MM-ddTHH:mm:ss")
// const fromTime = ref(from)
// const toTime = ref(to)
const filterTxt = ref('')

const log = ref('')

function selTime(event, value) {
  console.log(event)
  timeArr.forEach(t => {
    if (t.value == value) {
      console.log('set', t.sel, true)
      t.sel = true
    } else {
      t.sel = false
    }
  })
  minutesAgo.value = value
  console.log(minutesAgo.value, timeArr)
}

function query() {
  const deltaMin = minutesAgo.value
  const ft = new Date(Date.now() - deltaMin * 60 * 1000).Format("yyyy-MM-ddTHH:mm:ss")
  const tt = new Date(Date.now() + 60 * 1000).Format("yyyy-MM-ddTHH:mm:ss")

  Toast.loading({
    message: 'loading...',
    forbidClick: true,
    loadingType: 'spinner',
    duration: 10 * 1000
  });
  fetch(`/api/v1beta/miner/log?since=${ft}&until=${tt}&filter=${filterTxt.value}`)
    .then(r => r.json())
    .then(r => {
      if (r.code == 200) {
        log.value = r.data + '\n\n'
      } else {
        Dialog.alert({ message: "query log error:" + r.message })
      }
    }).catch(err => {
      Dialog.alert({ message: "query log error:" + err })
    }).finally(() => {
      Toast.clear()
    })
}

function fullScreen() {
  document.querySelector('#log-con').requestFullscreen()
}

</script>

<style lang="less" scoped>
.tag-tip {
  cursor: pointer;
}

.tag-sel {
  background-color: var(--van-tag-primary-color) !important;
  color: #fff !important;
}
.log {
  padding: 10px;
  background-color: #333;
  color: #fafafa;
  min-height: 50vh;
  line-height: 20px;
  font-size: 13px;
  overflow-y: scroll;
}
</style>