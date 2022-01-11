<template>
  <Cell>
    <div class="big-title">Control</div>
  </Cell>

  <CellGroup title="Device">
    <Cell title="Reboot Device">
      <Button size="small" type="danger" plain @click="reboot">Reboot</Button>
    </Cell>
    <Cell title="Resync Miner">
      <Button size="small" type="primary" plain @click="resync">Resync</Button>
    </Cell>
    <Cell title="Restart Miner">
      <Button size="small" type="primary" plain @click="restartMiner">Restart</Button>
    </Cell>
  </CellGroup>

  <CellGroup title="Snapshot">
    <Cell title="Generate Snapshot">
      <Button size="small" type="primary" plain @click="snapshot">Generate</Button>
    </Cell>
    <Cell :title="'Snapshot File :' + snapState.time">
      <Button
        v-if="snapState.state == 'done'"
        size="small"
        type="primary"
        plain
        @click="download"
      >Download</Button>
    </Cell>
    <Cell title="Apply Snapshot">
      <input id="file" class="hidden" ref="file" type="file" @change="handleFileChange" />
      <Button size="small" type="primary" plain @click="uploadSnapshot">Upload</Button>
    </Cell>
    <Cell>
      <Progress v-if="showProgress" :percentage="progress" :show-pivot="false"></Progress>
    </Cell>
  </CellGroup>
  <CellGroup title="Advanced">
    <Cell title="Workspace Update">
      <Button
        v-if="toUpdateWorkspace === true"
        size="small"
        type="danger"
        plain
        @click="updateWorkspace"
      >Update</Button>
      <Tag v-if="toUpdateWorkspace === false" type="success">up to date</Tag>
    </Cell>
    <Cell title="Workspace Reset">
      <Button size="small" type="danger" plain @click="resetWorkspace">Reset</Button>
    </Cell>
  </CellGroup>
</template>

<script setup>
import { reactive, ref } from "vue"
import { CellGroup, Cell, Button, Toast, Dialog, Progress, Notify, Tag } from 'vant'
import * as axios from "axios"
import * as api from "../api/backend"

const toUpdateWorkspace = ref(null)
const file = ref(null)
const showProgress = ref(false)
const progress = ref(0)
const snapState = reactive({ state: "unknown", file: "", time: "not generated" })

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
    message: 'processing ...',
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


function restartMiner() {
  Toast.loading({
    message: 'processing ...',
    forbidClick: true,
    loadingType: 'spinner',
    duration: 0
  });
  fetch('/api/v1/miner/restart', { method: 'POST' })
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

// open(`/api/v1/miner/snapshot/file/${r.data.file}`, "_blank")

function snapshot() {
  fetch("/api/v1/miner/snapshot", {
    method: "POST"
  }).then(r => r.json())
    .then(r => {
      console.log("snapshot res:", r)
      Dialog.alert({ message: "request success, refresh to check result after severial minuts" })
    })
    .catch(e => {
      console.error("snapshot error", e)
      Dialog.alert({ message: "snapshot error, try after several minutes" })
    })
}


function snapshotState() {
  fetch("/api/v1/miner/snapshot/state")
    .then(r => r.json())
    .then(r => {
      console.log("snapshot state:", r)
      if (r.code == 200) {
        if (r.data?.file && r.data?.state == 'done') {
          snapState.file = r.data.file
          snapState.time = r.data.time
          snapState.state = r.data.state
        }
      } else {
        Dialog.alert({ message: "load snapshot state error, please retry after a while" })
      }
    })
    .catch(e => {
      console.error("get snapshot state error", e)
      Dialog.alert({ message: "error, please retry after a while" })
    })
}

function download() {
  open(`/api/v1/miner/snapshot/file/${snapState.file}`, "_blank")
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
        Dialog.alert({ message: "request success, check result after severial minuts" })
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

function resetWorkspace() {
  Dialog.confirm({
    title: "Warning：Please operate carefully",
    message: "The device program and related configuration will be reset"
  }).then(() => {
    doResetWorkspace()
  }).catch(() => {
    // ignore
  })
}

async function checkWorkspaceUpdate() {
  api.checkWorkspaceUpdate()
    .then((r) => {
      toUpdateWorkspace.value = r
    })
    .catch(err => {
      const msg = err.response?.data?.message ? err.response.data.message : err.message
      Notify("failed to check workspace update:" + msg)
    })
}

function updateWorkspace() {
  api.workspaceUpdate()
    .then(r => {
      Dialog.alert({ message: "Request success. check after several minutes" });
    })
    .catch(err => {
      console.log('========', err.response?.data?.message)
      const msg = err.response?.data?.message ? err.response.data.message : err.message
      Dialog.alert({ type: "warning", message: msg })
    })
}

function doResetWorkspace() {
  fetch("/api/v1/workspace/reset", {
    method: "POST"
  }).then(r => r.json())
    .then(r => {
      if (r.status != 200) {
        console.log(r.status, r.data)
        Dialog.alert({ message: "error, http status " + r.status + "\n" + r.data?.message })
      } else {
        if (r.data.code != 200) {
          Dialog.alert({ message: "error " + r.data.message })
        } else {
          Dialog.alert({ message: "request success, check result after severial minuts" })
        }
      }

      console.log("workspace reset res:", r)
      Dialog.alert({ message: "success" })
    })
    .catch(e => {
      console.error("workspace reset error", e)
      Dialog.alert({ message: "snapshot error, try after several minutes" })
    })
}

snapshotState()
checkWorkspaceUpdate()

</script>

