<template>
  <CellGroup>
    <Cell title="Log Type">
      <RadioGroup v-model="logType" direction="horizontal">
        <Radio name="pktfwdLog">Packet Forward</Radio>
        <Radio name="minerLog">Miner</Radio>
      </RadioGroup> 
    </Cell>
    <Cell title="Recent Time" v-if="logType=='pktfwdLog'">
      <RadioGroup v-model="minutesAgo" direction="horizontal">
        <Radio v-for="item in timeArr"
          :name="item.value"
        >{{ item.label }}</Radio>
      </RadioGroup>
    </Cell>
    <!-- <Cell title="From Time">
      <input v-model="fromTime" type="datetime-local" />
    </Cell>
    <Cell title="To Time">
      <input v-model="toTime" type="datetime-local" />
    </Cell>-->
    <Cell title>
      <template #title>
        Filter Log&nbsp;
        <template v-if="logType=='pktfwdLog'">
          <Tag type="primary" plain round @click="filterTxt = 'JSON up'" class="tag-tip">Up link</Tag>&nbsp;
          <Tag type="primary" plain round @click="filterTxt = 'JSON down'" class="tag-tip">Down link</Tag>
        </template>
      </template>
      <Field v-model="filterTxt"
        placeholder="input filter text or select filter on the left" />
    </Cell>
    <Cell>
      <template #title>
        <span  v-if="logType!='pktfwdLog'">
          Show up to {{limitLine}} lines
        </span>
      </template>
      <Button type="primary" size="small" plain @click="fullScreen">Full Screen Log</Button>&nbsp;
      <Button type="primary" size="small" @click="query">Query</Button>
    </Cell>
  </CellGroup>

  <!-- <textarea v-model="log" style="width:100%"></textarea> -->
  <!-- <Divider /> -->
  <!-- <Field v-model="log" type="textarea" autosize disabled class="log"></Field> -->
  <pre id="log-con" class="log">
    <div v-for="l in logs" class="log-msg">
      <span class="log-date">{{l.time}}</span>
      {{l.message}}
    </div>
  </pre>
  <div class="van-safe-area-bottom"></div>
  <br/>
</template>

<script setup>
import { ref, reactive } from "vue"
import { CellGroup, Cell, Button, Field, Dialog, Divider, Tag, Toast, RadioGroup, Radio } from "vant"
import * as api from "../api"

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
const logType = ref('pktfwdLog')
const filterTxt = ref('')

const logs = reactive([])
const limitLine = ref(1000)

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

  const params = {
    logType: logType.value,
    filter: filterTxt.value,
    fromTime: ft,
    toTime: tt,
    limitLine: limitLine.value
  }
  api.logQuery(params)
    .then(r => {
      const arr = r.split('\n')
      logs.splice(0, logs.length)
      for (const l of arr) {
        try {
          if (l !== '') {
            const lj = JSON.parse(l)
            const t = Number(lj.time)
            if (isNaN(t)) {
              lj.time = lj.time.substr(5, 14)
            } else {
              lj.time = new Date(t/1000).Format("MM-dd HH:mm:ss")
            }
            logs.push(lj)
          }
        } catch (e) {
          console.error('parse log error: ', e)
        }
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

.van-field {
  border: 1px solid #efefef;
  border-radius: 100px;
}

.log {
  padding: 0px 10px;
  background-color: #454545;
  height: calc(100vh - 270px);
  color: #fafafa;
  line-height: 20px;
  font-size: 14px;
  overflow-y: scroll;
  .log-msg {
    white-space: nowrap;
    line-height: initial;
    color: #27aa5e;
  }
  .log-date {
    color: #8b959c;
    // user-select: none;
  }
}
</style>
