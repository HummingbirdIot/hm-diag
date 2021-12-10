<template>
  <Cell>
    <div class="big-title">Control</div>
  </Cell>

  <CellGroup>
    <Cell title="Reboot device">
      <Button size="small" type="danger" plain @click="reboot">reboot</Button>
    </Cell>
    <Cell title="Resync miner">
      <Button size="small" type="primary" plain @click="resync">Resync</Button>
    </Cell>
  </CellGroup>

  <CellGroup>
    <Cell title="Generate snapshot">
      <Button size="small" type="primary" plain @click="snapshot">Download</Button>
    </Cell>
    <Cell title="Upload snapshot">
      <Button size="small" type="primary" plain  @click="uploadSnapshot">Upload</Button>
    </Cell>
  </CellGroup>
</template>

<script setup>
import { CellGroup, Cell, Button, Toast, Dialog } from 'vant'

function reboot() {
  fetch('/api/v1/device/reboot', { method: 'POST' })
    .then(r => r.json())
    .then(r => {
      if (r && r.code == 200) {
        Toast.success("success")
      } else {
        Dialog.alert({ message: "error:" + r.message })
      }
    })
    .catch(r => {
      Dialog.alert({ message: "error:" + r })
    })
}

function resync() {
  Toast.loading({
    message: '加载中...',
    forbidClick: true,
    loadingType: 'spinner',
    duration: 0
  });
  fetch('/api/v1/miner/resync', { method: 'POST' })
    .then(r => r.json())
    .then(r => {
      Toast.clear()
      if (r && r.code == 200) {
        Toast.success("success")
      } else {
        Dialog.alert({ message: "error:" + r.message })
      }
    })
    .catch(r => {
      Dialog.alert({ message: "error:" + r })
    })
}

function snapshot() {
  alert("todo")
}

function uploadSnapshot() {
  alert("todo")
}

</script>

