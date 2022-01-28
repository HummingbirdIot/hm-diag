<template>
  <Cell>
    <div class="big-title">Setting</div>
  </Cell>

  <CellGroup title="Github repository proxy">
    <Cell>
      <Field v-model="rp" placeholder="mirror url, support http:// and https://">
        <template #button>
          <Button size="small" type="primary" @click="confirm('gitRepo')">Confirm</Button>
        </template>
      </Field>
    </Cell>
  </CellGroup>

  <CellGroup title="Github release file proxy">
    <Cell>
      <Field v-model="rrp" placeholder="mirror url, support http:// and https://">
        <template #button>
          <Button size="small" type="primary" @click="confirm('gitRelease')">Confirm</Button>
        </template>
      </Field>
    </Cell>
  </CellGroup>
  <br />
  <br />
  <br />
</template>

<script setup>
import { ref } from 'vue'
import { CellGroup, Cell, Field, Button, Dialog, Toast } from 'vant'
import * as api from '../api'
import * as errors from '../util/errors'


const rp = ref("")
const rrp = ref("")

api.proxyConfigGet("gitRepo")
  .then(rpRe => {
    rp.value = rpRe?.value ? rpRe.value : ""
    console.log(rp.value)
  })
api.proxyConfigGet("gitRelease")
  .then(rrpRe => {
    rrp.value = rrpRe?.value ? rrpRe.value : ""
    console.log(rrp.value)
  })

function confirm(item) {
  const v = item == 'gitRepo' ? rp.value : rrp.value
  let u = v
  if (u != "") {
    u = getUrl(v)
  }
  const type = item == 'gitRepo' ? 'mirror' : 'urlPrefix';
  api.proxyConfigSet(item, {type, value: u})
    .then(r => {
        Toast.success("set success")
    })
    .catch(err => {
        Dialog.alert({ message: "set error :" + errors.getMsg(err)})
    })
}

function getUrl(v) {
  try {
    const url = new URL(v)
    if (url.protocol !== 'http:' && url.protocol !== 'https:') {
      Dialog.alert({ message: "wrong url, scheme should be http or https" })
      return null
    }
    return url.origin + '/'
  } catch (err) {
    Dialog.alert({ message: "set error :" + r?.message })
    return null
  }
}

</script>

