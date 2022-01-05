<template>
  <div class="page">
    <!-- <img src="./asset/logo-long.png" alt class="logo" /> -->
    <Row justify="center">
      <h1>Hotpot Info</h1>
    </Row>
    <Row justify="center">
      <h2>{{ data?.miner?.infoSummary?.name }}</h2>
    </Row>

    <CellGroup>
      <Cell title="Helium Address" @click="openHeliumExplorer">
        <!-- <Button :url="heliumAddr" plain type="default" size="small">View on Explore</Button> -->
        <a href="#">View on Explore</a>
      </Cell>
      <Cell title="Height Status">{{ data?.miner?.infoSummary?.height }} / {{ heliumHeight }}</Cell>
      <Cell title="Miner Version">{{ data?.miner?.infoSummary?.version }}</Cell>
      <Cell title="Firmware Version">{{ data?.miner?.infoSummary?.firmwareVersion?.split('\n').slice(-2, -1)[0].split('=')[1].replaceAll('"', '') }}</Cell>
      <Cell title="Region Plan">{{ region }}</Cell>
      <Cell title="Miner Connected to Blockchain">
        <Tag
          v-if="data?.miner?.peerBook?.length > 0 && data?.miner?.peerBook[0]?.connectionCount > 0"
          type="success"
        >True</Tag>
        <Tag v-else type="warning">Flase</Tag>
      </Cell>
      <Cell title="Miner Relayed">
        <Tag v-if="data?.miner?.infoP2pStatus?.natType == 'symmetric'" type="warning">True</Tag>
        <Tag v-else type="success">Flase</Tag>
      </Cell>
    </CellGroup>

    <CellGroup title=" ">
      <Cell title="CPU Temperature">{{ data?.device?.cpuTemp }}</Cell>
      <Cell title="Disk Used">
        <Progress class="mg-top" :percentage="diskPercentage" :show-pivot="true" :stroke-width="4"></Progress>
      </Cell>
      <Cell title="Memory Used">
        <Progress class="mg-top" :percentage="memPercentage" :show-pivot="true" :stroke-width="4"></Progress>
      </Cell>
      <Cell
        title="ETH0 MAC"
      >{{ data?.device?.netInterface?.find(i => i.name == 'eth0')?.hardwareAddr }}</Cell>
      <Cell
        title="WLAN0 MAC"
      >{{ data?.device?.netInterface?.find(i => i.name == 'wlan0')?.hardwareAddr }}</Cell>

      <Cell title="Miner Log" is-link to="/minerLog">
      </Cell>
    </CellGroup>
    <br/>
    <br/>
    <br/>
  </div>
  
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue';
import {
  Row,
  Col,
  Tag,
  Icon,
  Cell,
  CellGroup,
  Toast,
  Button,
  Progress,
  Divider,
} from 'vant';

const data = reactive({})
const heliumAddr = ref('')
const heliumHeight = ref()
const diskPercentage = ref(0)
const memPercentage = ref(0)
const region = ref('')

function openHeliumExplorer() {
  window.open(heliumAddr.value, '__blank')
}

function fetchData() {
  Toast.loading({
    message: 'loading...',
    forbidClick: true,
    loadingType: 'spinner',
    duration: 10 * 1000
  });
  fetch('/state')
    .then(r => r.json())
    .then(r => {
      Toast.clear()
      Object.assign(data, r.data)
      fillData(data)
    })
    .catch(r => {
      Toast.fail("error :" + r)
    })
}

function fillData(data) {
  region.value = data?.miner?.infoRegion ? data.miner.infoRegion.split('_')[1] : ''
  region.value = region.value.toUpperCase()
  heliumAddr.value = `https://explorer.helium.com/hotspots/${data?.miner?.peerAddr.slice(5)}`
  diskPercentage.value = data?.device?.disk[0].usedPercent ? Math.ceil(data.device.disk[0].usedPercent) : 0
  memPercentage.value = Math.ceil(data?.device?.mem ? data.device.mem.used * 100 / data.device.mem.total : 0)
}

function fetchHeliumHeight() {
  const api = 'https://api.helium.io/v1/blocks/height'
  fetch(api)
    .then(r => r.json())
    .then(r => {
      heliumHeight.value = r.data.height
    })
}

onMounted(() => {
  fetchHeliumHeight()
  fetchData();
})

</script>

<style lang="less" scoped>
.logo {
  // max-width:100px;
  max-width: 80%;
  display: block;
  margin: 0 auto;
}
.mg-top {
  margin-top: 10px;
}
</style>
