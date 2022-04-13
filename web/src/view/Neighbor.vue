<template>
  <CellGroup>
    <Cell>
      <div class="big-title">
        Neighbors
        <Icon name="question-o" @click="showTip"></Icon>
      </div>
    </Cell>
    <Search v-model="searchTxt" placeholder="search hotspot name, version, ip"></Search>
    <Cell>
      <div style>
        <Icon name="info-o"></Icon>
        <span class="sum-info">
          Total: {{
            filteredNeighbors?.length != undefined
              ? filteredNeighbors.length
              : '-'
          }}
        </span>
        <span class="sum-info">
          Syncing: {{
            heliumHeight != 0
              ? filteredNeighbors?.filter(n => heliumHeight - n.height > 1).length
              : '-'
          }}
        </span>
        <span class="sum-info">
          Relayed: {{
            filteredNeighbors?.filter(n => n?.isRelay)?.length
          }}
        </span>
      </div>
    </Cell>
    <Cell v-for="n in filteredNeighbors" is-link @click="openHotspotWeb(n)">
      <template #title>
        <div class="pot-name">{{ n.showName ? n.showName : n.address }}</div>
        <template v-if="n.dataReady">
          <div class="pot-info">IP: <span class="pot-value">{{ n.address }}</span></div>
          <div class="pot-info">Block Height: <span class="pot-value">{{ n.height }} / {{ heliumHeight }}</span></div>
          <div class="pot-info">Miner Version: <span class="pot-value">{{ n.version }}</span></div>
          <div class="pot-info">Listen Address:
            <span class="pot-value">{{n.listenAddress}}</span>
          </div>
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
  <br />
  <br />
  <br />
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { CellGroup, Cell, Search, Dialog, Toast, Tag, Icon, Notify } from 'vant'
import * as api from '../api'
import * as errors from '../util/errors'

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
  api.lanHotspots()
    .then(r => {
      neighbors.value = r
      filteredNeighbors.value = neighbors.value
      fetchHotspotsInfo()
    })
    .catch(err => {
      Dialog.alert({ message: "get hotspots list error :" + errors.getMsg(err) })
    })
    .finally(() => {
      Toast.clear()
    })
}

function openHotspotWeb(d) {
  const addr = `http://${d.address}:${d.port}`
  window.open(addr, '_blank')
}

function fetchHotspotsInfo() {
  for (const n of neighbors.value) {
    if (!n.address || !n.port) {
      console.error('no address or port', n)
      continue
    }
	  // TODO change this route to /state after next version
    const api = `http://${n.address}:${n.port}/state?cache=true`
    const ctrl = new AbortController()
    const timeoutId = setTimeout(() => ctrl.abort(), 5000)
    fetch(api, { signal: ctrl.signal })
      .then(r => r.json())
      .then(r => {
        if (r.code == 200) {
          console.log('go hotspot info', r.data)
          const s = r.data.miner.infoSummary
          const p2p = r.data.miner.infoP2pStatus
          const pb = r.data.miner.peerBook
          n.showName = s.name
          n.isRelay = p2p.natType == 'symmetric' ? true : false
          n.height = s.height
          n.version = s.version
          n.dataReady = true
          n.listenAddress = pb.length > 0 ? pb[0].listenAddresses.length > 0 ? pb[0].listenAddresses[0] : '' : ''
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
  //2022.4.13，因为helium api的调用限制，先将height的显示都暂时取消
  api.blockHeight().then(h => {
    heliumHeight.value = h
  }).catch(err => {
    Dialog.alert({ message: "Failed to load helium block height" });
    console.log("Failed to load helium block height", err);
  })
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
.sum-info {
  margin-left: 6px;
  margin-right: 6px;
  font-size: 14px;
  font-weight: 300;
  color: #333;
}
.pot-name {
  font-size: 16px;
  font-weight: 500;
}
.pot-info {
  color: #666;
  font-size: 12px;
  margin-right: -90%;
  .pot-value {
    font-weight: 500;
  }
}
</style>

