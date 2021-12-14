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
      <input id="file" class="hidden" ref="file" type="file" @change="handleFileChange" />
      <Button size="small" type="primary" plain @click="uploadSnapshot">Upload</Button>
    </Cell>
    <Cell>
      <Progress v-if="showProgress" :percentage="progress" :show-pivot="false"></Progress>
    </Cell>
  </CellGroup>
</template>

<script setup>
import { ref } from "vue"
import { CellGroup, Cell, Button, Toast, Dialog, Progress, Notify } from 'vant'
import * as axios from "axios"

const file = ref(null)
const showProgress = ref(false)
const progress = ref(0)

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
  fetch("/api/v1/miner/snapshot", {
    method: "POST"
  }).then(r => r.json())
    .then(r => {
      console.log("snapshot res:", r)
      Dialog.alert({ message: "success snapshot" })
      open(`/api/v1/miner/snapshot/file/${r.data.file}`, "_blank")
    })
    .catch(e => {
      console.error("snapshot error", e)
      Dialog.alert({ message: "snapshot error, try after several minutes" })
    })
}

function handleFileChange() {
  console.log("seleted file", file.value.files)
  const fd = new FormData()
  fd.append("file", file.value.files[0])
  showProgress.value = true
  Notify({ type: "warning", duration: 5 * 1000, message: "Don't operate till the end" })
  axios.default.post("/api/v1/miner/snapshot/apply", fd, {
    onUploadProgress: (e) => {
      progress.value = Math.round(e.loaded * 100 / e.total)
      if (progress.value == 100) {
        Notify({ type: "success", message: "Upload done, waite apply it" })
        Toast.loading({ duration: 0, message: "applying snapshot" })
      }
    }
  }).then(r => {
    if (r.status != 200) {
      console.log(r.status, r.data)
      Dialog.alert({ message: "error, http status " + r.status + "\n" + r.data?.message })
    } else {
      if (r.data.code != 200) {
        Dialog.alert({ message: "error " + r.data.message })
      } else {
        Dialog.alert({ message: "success" })
      }
    }
  })
    .catch(e => {
      console.log(e)
      Dialog.alert({ message: "" + e + "\n" + e.response?.data?.message })
    })
    .finally(() => {
      showProgress.value = false
      Toast.clear()
    })

}
function uploadSnapshot() {
  document.querySelector("#file").click()
}

</script>

