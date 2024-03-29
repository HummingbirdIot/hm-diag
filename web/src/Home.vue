<template>
  <div class="page">
    <!-- <img src="./asset/logo-long.png" alt class="logo" /> -->
    <Row justify="center">
      <h1>Hotspot Info</h1>
    </Row>
    <Row justify="center">
      <h2>{{ data?.miner?.infoSummary?.name }}</h2>
    </Row>

    <CellGroup>
      <Cell title="Helium Address" @click="openHeliumExplorer">
        <a v-if="heliumAddr != ''" href="#">View on Explore</a>
      </Cell>
      <Cell>
        <template #title>
          Height Status
          <Icon name="question-o" @click="blockHeightTip"></Icon>
        </template>
        {{ data?.miner?.infoSummary?.height }} / {{ heliumHeight }}
      </Cell>
      <Cell title="Miner Version">{{ data?.miner?.infoSummary?.version }}</Cell>
      <Cell title="Firmware Version">{{ firmwareVersion }}</Cell>
      <Cell title="Region Plan">{{ region }}</Cell>
      <Cell title="Miner Connected to Blockchain" :label=data?.miner?.gatewayInfo?.gateway.uri>
        <Tag
          v-if="data?.miner?.peerBook?.length > 0 && data?.miner?.peerBook[0]?.connectionCount > 0 || data?.miner?.gatewayInfo?.gateway.uri"
          type="success"
        >True</Tag>
        <Tag v-else type="warning">False</Tag>
      </Cell>
      <Cell title="Miner Relayed" v-if="data?.miner?.infoP2pStatus?.natType">
        <Tag v-if="data?.miner?.infoP2pStatus?.natType == 'symmetric'" type="warning">True</Tag>
        <Tag v-else type="success">False</Tag>
      </Cell>
      <Cell title="Miner Listen Address" v-if="data?.miner?.peerBook">
        {{ data?.miner?.peerBook?.length > 0 
          ? data.miner.peerBook[0].listenAddresses?.length > 0
            ? data.miner.peerBook[0].listenAddresses[0] : ''
          : ''
        }}
      </Cell>
      <Cell title="Version">
        {{version?.version}}
      </Cell>
    </CellGroup>

    <CellGroup title=" ">
      <Cell title="CPU Temperature">{{ data?.device?.cpuTemp }}</Cell>
      <Cell title="CPU Frequency">
        {{
          data?.device?.cpuFreq
            ? Math.round(data.device.cpuFreq / 1e5) / 10 + 'GHz'
            : ''
        }}
      </Cell>
      <Cell>
        <template #title>
          CPU Percentage
          <Icon name="question-o" @click="cpuPercentTip"></Icon>
        </template>
        <template #value>
          <span v-for="item in data?.device?.cpuPercent" :key="item">{{ Math.round(item) }}%&nbsp;</span>
        </template>
        <!-- <template #value>
            <Progress
              v-for="(item,i) in data?.device?.cpuPercent"
              :percentage="Math.round(item * 100) / 100"
              stroke-width="2"
              style="margin: 18px 0px"
            ></Progress>
        </template>-->
      </Cell>
      <Cell title="Disk Used">
        <Progress class="mg-top" :percentage="diskPercentage" :show-pivot="true" :stroke-width="2"></Progress>
      </Cell>
      <Cell title="Memory Used">
        <Progress class="mg-top" :percentage="memPercentage" :show-pivot="true" :stroke-width="2"></Progress>
      </Cell>
      <Cell
        title="ETH0 MAC"
      >{{ data?.device?.netInterface?.find(i => i.name == 'eth0')?.hardwareAddr }}</Cell>
      <Cell
        title="WLAN0 MAC"
      >{{ data?.device?.netInterface?.find(i => i.name == 'wlan0')?.hardwareAddr }}</Cell>

    </CellGroup>

    <CellGroup title=" ">
      <Cell title="Network Test" is-link to="/networkTest"></Cell>
      <Cell title="Log Query" is-link to="/logQuery"></Cell>
      <Cell title="Download Log">
        <Button size="small" type="primary" plain @click="downloadAllLog">Download</Button>
      </Cell>
    </CellGroup>
    <br />
    <br />
    <br />
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue';
import {
  Row,
  Tag,
  Cell,
  CellGroup,
  Toast,
  Progress,
  Dialog,
  Icon,
  Button,
  Notify
} from 'vant';
import * as api from './api'
import * as errors from './util/errors'

const data = reactive({})
const heliumAddr = ref('')
const heliumHeight = ref()
const diskPercentage = ref(0)
const memPercentage = ref(0)
const region = ref('')
const firmwareVersion = ref('')
const version = reactive({})

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
  api.stateGet()
    .then(r => {
      Toast.clear()
      localStorage.setItem("macPath",r?.device?.netInterface?.find(i => i.name == 'eth0')?.hardwareAddr)
      Object.assign(data, r)
      fillData(data)
    })
    .catch(err => {
      console.error(err)
      Dialog.alert({ message: "get state error:" + errors.getMsg(err) })
    }).finally(() => {
      checkHotspotOnboarding()
    })
}

function fillData(data) {
  const { miner, device } = data
  region.value = miner?.infoRegion
  heliumAddr.value = `https://explorer.helium.com/hotspots/${miner?.peerAddr?.slice(5)}`
  diskPercentage.value = device?.disk?.length > 0 ? Math.ceil(device.disk[0].usedPercent) : 0
  memPercentage.value = device?.mem ? Math.ceil(device.mem.used * 100 / device.mem.total) : 0
  firmwareVersion.value = miner?.infoSummary?.firmwareVersion
    ?.split('\n')
    ?.slice(-2, -1)[0]
    ?.split('=')[1]
    ?.replaceAll('"', '')
}

function downloadAllLog(){
  Notify({ type: 'primary', message: 'Downloading, please wait patiently' });
  api.downloadLog().then(res => {
    let fileName = res.headers["content-disposition"].split("url=")[1];
    console.log(fileName)
    const blob = new Blob([res.data],{type: 'application/x-tar'});
    const elink = document.createElement('a');
    elink.download = fileName;
    elink.style.display = 'none';
    elink.href = URL.createObjectURL(blob);
    document.body.appendChild(elink);
    elink.click();
    URL.revokeObjectURL(elink.href);// 释放URL 对象
    document.body.removeChild(elink);
    Notify({ type: "success", message: "Download log successful" })
  }).catch(err=>{
    console.error(err)
    Notify({type: 'danger', message:"failed to download log:" + errors.getMsg(err)})
  })
}

function fetchBlockHeight() {
  api.blockHeight()
  .then(r => {
    heliumHeight.value = r
  }).catch(e => {
    api.blockHeightFromHelium()
      .then(r => {
        heliumHeight.value = r
      }).catch(err => {
        Dialog.alert({ message: "Failed to load helium block height" });
        console.log("Failed to load helium block height", err);
      })
  })
}

function blockHeightTip() {
  Dialog.alert({
    message: 'Block height of Miner / Block height from Helium API.'
      + ' Helium API data may be a delay of some time'
  })
}
function cpuPercentTip() {
  Dialog.alert({ message: 'Percentage per CPU' })
}

function checkHotspotOnboarding(){
  api.checkOnboarding().then(res => {
    if(!res && !data?.miner?.infoSummary?.owner){
      Dialog.alert({ title:"Warning",message: "Hotspot is not onboarding" });
    }
  }).catch(err =>{
    Dialog.alert({ message: "Failed to load hotspot onboarding" });
  })
}

onMounted(() => {
  api.version().then(r=>Object.assign(version, r))
  fetchBlockHeight()
  fetchData()
  // const mock = { "data": { "device": { "cpuFreq": "1500398464", "cpuPercent": [14.99999999996362, 27.638190954866655, 22.110552764180127, 42.21105527640396], "cpuTemp": "46.7'C", "disk": [{ "free": 6101557248, "fstype": "ext2/ext3", "path": "/", "total": 30495752192, "used": 23427592192, "usedPercent": 79.33717237471504 }], "host": { "hostname": "raspberrypi", "uptime": 152458, "bootTime": 1641465242, "procs": 152, "os": "linux", "platform": "raspbian", "platformFamily": "debian", "platformVersion": "10.11", "kernelVersion": "5.10.63-v8+", "kernelArch": "aarch64", "virtualizationSystem": "", "virtualizationRole": "", "hostId": "591d3131-3b85-437b-9537-72386ea4e881" }, "mem": { "available": 3520720896, "buffers": 43294720, "cached": 2211631104, "free": 1268211712, "shared": 462848, "total": 3979366400, "used": 456228864 }, "netInterface": [{ "index": 1, "mtu": 65536, "name": "lo", "hardwareAddr": "", "flags": ["up", "loopback"], "addrs": [{ "addr": "127.0.0.1/8" }, { "addr": "::1/128" }] }, { "index": 2, "mtu": 1500, "name": "eth0", "hardwareAddr": "e4:5f:01:30:54:27", "flags": ["up", "broadcast", "multicast"], "addrs": [{ "addr": "172.20.0.197/24" }, { "addr": "fe80::c2a:3b78:1e4e:f6bc/64" }] }, { "index": 3, "mtu": 1500, "name": "wlan0", "hardwareAddr": "e4:5f:01:30:54:2a", "flags": ["up", "broadcast", "multicast"], "addrs": [{ "addr": "192.168.88.192/21" }, { "addr": "fe80::9dce:8a57:7c5d:f172/64" }] }, { "index": 4, "mtu": 1500, "name": "docker0", "hardwareAddr": "02:42:dc:6c:47:50", "flags": ["up", "broadcast", "multicast"], "addrs": [{ "addr": "172.17.0.1/16" }] }], "wifi": { "connected": true, "name": "WiFi", "powered": true, "tethering": false, "type": "wifi" } }, "miner": { "infoHeight": 1171072, "infoP2pStatus": { "connected": "yes", "dialable": "yes", "height": 1171072, "natType": "none" }, "infoRegion": "CN470", "infoSummary": { "firmwareVersion": "DISTRIB_ID=hummingbird\nDISTRIB_RELEASE=2021.12.29.2\nDISTRIB_DESCRIPTION=\"Helium H500 Firmware, 2021.12.29.2\"\nDISTRIB_DATE=\"2022.1.3 15:09:29\"\n", "height": 1171072, "name": "howling-cedar-shark", "version": "2021.12.29.2" }, "peerAddr": "/p2p/11akXhBWtRsYemycJoHEeFkFQW2XrwuELqjQXCbR42cuo559sLn", "peerBook": [{ "address": "/p2p/11akXhBWtRsYemycJoHEeFkFQW2XrwuELqjQXCbR42cuo559sLn", "connectionCount": 8, "lastUpdated": "232.924", "listenAddrCount": 1, "listenAddresses": ["/ip4/47.241.29.36/tcp/46197"], "name": "howling-cedar-shark", "nat": "none", "sessions": [{ "local": "/ip4/172.20.0.197/tcp/46197", "name": "sharp-clay-deer", "p2p": "/p2p/1128xhvBgYH3Uiw33ujXHA2zaruaPisDwtKgJpgshN9w7C5XS8vv", "remote": "/ip4/178.198.74.198/tcp/44158" }, { "local": "/ip4/172.20.0.197/tcp/46197", "name": "straight-burgundy-lizard", "p2p": "/p2p/112BLWtDfyb5cqAXvf24EtNT3RpK3weCqHhjt6wUC7FuEH3rWFoK", "remote": "/ip4/2.81.105.112/tcp/44158" }, { "local": "/ip4/172.20.0.197/tcp/46197", "name": "alert-canvas-pelican", "p2p": "/p2p/112eJ85WeNGjn1tPvarPBmCBzi6Nuoa2Y2A3rCRQLF4HN4q8DyGD", "remote": "/ip4/116.225.96.46/tcp/44158" }, { "local": "/ip4/172.20.0.197/tcp/46197", "name": "melodic-mango-perch", "p2p": "/p2p/11C9v6rWqob6ihz4BX2bfXVw2Me2EwJXUw3HxWTgrWeZiUxjNyo", "remote": "/ip4/27.155.101.105/tcp/44158" }, { "local": "/ip4/172.20.0.197/tcp/46197", "name": "ancient-beige-hare", "p2p": "/p2p/11UCbsgywnxzfXdjDxsHcUqwfPgjt9nSt9m3Kt3zQH4Phcz7DwT", "remote": "/ip4/54.251.69.171/tcp/2154" }, { "local": "/ip4/172.20.0.197/tcp/46197", "name": "spicy-azure-capybara", "p2p": "/p2p/11jfoconYsAqs1JNNA2VY3nyVdbTPcXFubbGtcrAH8UbvAX6ngu", "remote": "/ip4/73.84.57.66/tcp/44158" }, { "local": "/ip4/172.20.0.197/tcp/46197", "name": "bouncy-tan-cobra", "p2p": "/p2p/11ne6mfbupmogaj98k5uQqVggCzRbp4oh6FuQeAkLzzR2798R5C", "remote": "/ip4/54.66.57.169/tcp/443" }] }] }, "notice": "do not use this api path \"/\" to integrate, use api under path \"api/\"" }, "code": 200, "message": "OK" }
  // Object.assign(data, mock.data)
  // fillData(data)
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
