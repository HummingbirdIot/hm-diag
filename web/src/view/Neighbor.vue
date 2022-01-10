<template>
  <CellGroup>
    <Cell>
      <div class="big-title">
        Neighbors
        <Icon name="question-o" @click="showTip"></Icon>
      </div>
    </Cell>
    <Search v-model="searchTxt" placeholder="search hotspot name, version, ip"></Search>
    <Cell v-for="n in filteredNeighbors" is-link @click="openHotspotWeb(n)">
      <template #title>
        <div class="pot-name">{{ n.showName ? n.showName : n.address }}</div>
        <template v-if="n.dataReady">
          <div class="pot-ip">IP: {{ n.address }}</div>
          <div class="pot-ip">Block Height: {{ n.height }} / {{ heliumHeight }}</div>
          <div class="pot-ip">Miner Version: {{ n.version }}</div>
        </template>
      </template>
      <template #value>
        <Tag
          v-if="heliumHeight > 0 && n.height > 0"
          :type="heliumHeight - n.height <= 1 ? 'success' : 'warning'"
        >{{ heliumHeight - n.height <= 1 ? 'Synced' : 'Syncing' }}</Tag>&nbsp;
        <Tag
          v-if="n.isRelay !== undefined"
          :type="n.isRelay ? 'warning' : 'success'"
        >{{ n.isRelay ? 'Relayed' : 'NoRelay' }}</Tag>
      </template>
    </Cell>
  </CellGroup>
  <div class="van-safe-area-bottom"></div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from "vue"
import { CellGroup, Cell, Search, Dialog, Toast, Tag, Icon, Notify } from 'vant'

let searchTxt = ref('')
let neighbors = ref([])
let filteredNeighbors = ref([])
let heliumHeight = ref(0)

function showTip() {
  Dialog.alert({
    message: 'In order to speed up access, '
      + 'cached data is used, '
      + 'with a maximum delay of 30 seconds'
      + 'If the client and the device are not in the same local area network, '
      + 'detailed data cannot be seen'
  })
}

function loadLanHotspots() {
  Toast.loading({
    message: 'loading ...',
    forbidClick: true,
    loadingType: 'spinner',
    duration: 10 * 1000
  });
  fetch(`/api/v1/lan/hotspot`)
    .then(r => r.json())
    .then(r => {
      if (r && r.code == 200) {
        neighbors.value = r.data
        filteredNeighbors.value = neighbors.value
        fetchHotspotsInfo()
      } else {
        Dialog.alert({ message: "get hotspots list error :" + r?.message })
      }
    }).finally(() => {
      Toast.clear()
    })
}

function openHotspotWeb(d) {
  const addr = `http://${d.address}:${d.port}`
  window.open(addr, '_blank')
}

function fetchHeliumHeight() {
  const api = 'https://api.helium.io/v1/blocks/height'
  fetch(api)
    .then(r => r.json())
    .then(r => {
      heliumHeight.value = r.data.height
    }).catch(err => {
      Dialog.alert({ message: 'Failed to load helium block height' })
      console.log('Failed to load helium block height', err)
    })
}

function fetchHotspotsInfo() {
  for (const n of neighbors.value) {
    const api = `http://${n.address}:${n.port}/state?cache=true`
    const ctrl = new AbortController()
    const timeoutId = setTimeout(() => ctrl.abort(), 5000)
    fetch(api, { signal: ctrl.signal })
      .then(r => r.json())
      .then(r => {
        if (r.code == 200) {
          console.log('go spot info', r.data)
          const s = r.data.miner.infoSummary
          const p2p = r.data.miner.infoP2pStatus
          n.showName = s.name
          n.isRelay = p2p.natType == 'symmetric' ? true : false
          n.height = s.height
          n.version = s.version
          n.dataReady = true
        } else {
          n.dataReady = false
          n.showName = n.address
        }
      }).catch(e => {
        console.error('fetch hotspot info error', n, e)
        n.dataReady = false
        n.showName = n.address
        Notify({
          type: 'warning',
          message: 'Failed to get hotspot information, '
            + 'make sure you and the hotspot are in the same LAN'
        })
      })
      .finally(() => {
        clearTimeout(timeoutId)
      })
  }
}

watch(searchTxt, v => {
  filteredNeighbors.value = neighbors.value
    .filter(n =>
      v == ''
      || n.showName?.indexOf(v) >= 0
      || n.version?.indexOf(v) >= 0
      || n.address?.indexOf(v) >= 0
    )
})

onMounted(async () => {
  fetchHeliumHeight()
  loadLanHotspots()
  // const mockData = [
  //   {
  //     "name": "boxy-silver-orca-1712",
  //     "host": "boxy-silver-orca-1712.local",
  //     "address": "192.168.89.45",
  //     "port": 80
  //   }
  // ]
  // data.value = JSON.stringify(mockData, null, 2)
  // // Object.assign(neighbors, mockData)
  // mockData.forEach((item) => {
  //   neighbors.push(item)
  // })

})

</script>

<style lang="less" scoped>
.pot-name {
  font-size: 16px;
  font-weight: 500;
}
.pot-ip {
  color: #666;
  font-size: 12px;
  margin-right: -90%;
}
</style>

