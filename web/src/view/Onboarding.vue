<template>
  <Cell>
    <div class="big-title">Onboarding</div>
  </Cell>
  <CellGroup>
      <Field v-model="ownerAddress" label="Owner"
        placeholder="input owner address" >
      <template #button>
        <Button size="small" type="primary" plain @click="genTxn">Confirm</Button>
      </template>
    </Field>
    <div class="qr-con">
      <canvas id="qr"></canvas>
    </div>
    <div v-if="showQrTip" class="qr-tip">
      <p>
        Scan QR image above by <strong>Hummingbird Maker APP</strong> or Helium APP to continue onboarding.
      </p>
      <p>Or click <strong><a :href="qrLink">Here</a></strong> on your cellphone to call Helium APP for onboarding.</p>

    </div>
  </CellGroup>

</template>

<script setup>
import { ref } from "vue"
import { CellGroup, Cell, Button, Field, Dialog, Notify, Toast, } from "vant"
import * as api from "../api"
import Qrious from "qrious"

const ownerAddress = ref('')
const showQrTip = ref(false)
const qrLink = ref('')

function genTxn() {
  if(ownerAddress.value.length == 0 || ownerAddress.value.length != 51) {
    Notify('Please input valid address')
    return
  }
  Toast.loading({
    message: 'loading...',
    forbidClick: true,
    loadingType: 'spinner',
    duration: 10 * 1000
  });
  api.onboarding(ownerAddress.value)
    .then(txn=>{
      qrLink.value = `helium://add_gateway/:${txn}?from=hummingbird`
      console.log('txn qr data: ', qrLink.value)
      new Qrious({
        element: document.getElementById('qr'),
        value: qrLink.value,
        size: 300,
        level: 'L'
      });
      showQrTip.value = true
    })
    .catch(e=>{
      Dialog.alert({ message: 'error:' + e})
    })
    .finally(()=>{
      Toast.clear()
    })
}

</script>

<style lang="less" scoped>
.qr-con {
  display: flex;
  justify-content: center;
}
.qr-tip {
  text-align: center;
  padding: 10px;  
}
</style>
