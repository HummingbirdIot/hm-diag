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
  <br />
  <br />
  <br />
</template>

<script setup>
import { reactive, ref } from "vue"
import { CellGroup, Cell, Button, Toast, Dialog, Progress, Notify, Tag } from 'vant'
import * as axios from "axios"
import * as api from "../api/backend"
import * as errors from "../util/errors"

const toUpdateWorkspace = ref(null)
const file = ref(null)
const showProgress = ref(false)
const progress = ref(0)
const snapState = reactive({ state: "unknown", file: "", time: "not generated" })

function reboot() {
  api.deviceReboot()
    .then(r => {
      Toast.success("success")
    })
    .catch(err => {
      Dialog.alert({ message: "error:" + errors.getMsg(err) })
    })
}

function resync() {
  Toast.loading({
    message: 'processing ...',
    forbidClick: true,
    loadingType: 'spinner',
    duration: 0
  });
  api.minerResync() 
    .then(r => {
      Toast.clear()
      Toast.success("success")
    })
    .catch(err => {
      Dialog.alert({ message: "error:" + errors.getMsg(err) })
    })
}


function restartMiner() {
  Toast.loading({
    message: 'processing ...',
    forbidClick: true,
    loadingType: 'spinner',
    duration: 0
  });
  api.minerRestart() 
    .then(() => {
      Toast.clear()
      Toast.success("success")
    })
    .catch(err => {
      Dialog.alert({ message: "error:" + errors.getMsg(err) })
    })
}

function snapshot() {
  api.snap()
    .then(r => {
      Dialog.alert({ message: "request success, refresh to check result after severial minuts" })
    })
    .catch(e => {
      console.error("snapshot error", e)
      Dialog.alert({ message: "snapshot error, try after several minutes:\n" + errors.getMsg(e) })
    })
}


function snapshotState() {
  api.snapState()
  .then(r => {
    console.log("snapshot state:", r)
    if (r.file && r.state == 'done') {
      snapState.file = r.file
      snapState.time = r.time
      snapState.state = r.state
    }
  })
  .catch(e => {
    console.error("get snapshot state error", e)
    Dialog.alert({ message: "error, please retry after a while" })
  })
}

function download() {
  api.snapDownload(snapState.file)
}

function handleFileChange() {
  console.log("seleted file", file.value.files)
  const fd = new FormData()
  fd.append("file", file.value.files[0])
  showProgress.value = true
  Notify({ type: "warning", duration: 5 * 1000, message: "Don't operate till the end" })
  axios.default.post("/inner/api/v1/miner/snapshot/apply", fd, {
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

async function workspaceUpdateCheck() {
  api.workspaceUpdateCheck()
    .then((r) => {
      toUpdateWorkspace.value = r
    })
    .catch(err => {
      Notify("failed to check workspace update:" + errors.getMsg(err))
    })
}

function updateWorkspace() {
  api.workspaceUpdate()
    .then(r => {
      Dialog.alert({ message: "Request success. check after several minutes" });
    })
    .catch(err => {
      console.error('error:', err)
      Dialog.alert({ type: "warning", message: errors.getMsg(err) })
    })
}

function doResetWorkspace() {
  api.workspaceReset()
    .catch(err => {
      console.error("workspace reset error", err)
      Dialog.alert({ message: "snapshot error, try after several minutes : " + errors.getMsg(err) })
    })
}

snapshotState()
workspaceUpdateCheck()

</script>

