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
      <!-- <div class="qr-example">
        <img v-show="!showQr" src="../asset/img/qr.png" />
        <span></span>
      </div> -->
    </div>
    <div v-if="showQrTip" class="qr-tip">
        <p><strong>Scan</strong> QR image above by Hummingbird Maker APP or Helium APP to continue onboarding.</p>
        <!-- <p>Or <a :href="qrLink">click here on your cellphone</a> to call Helium APP for onboarding.</p> -->
        <p>The APP should have logged into the owner's account in advance.</p>
    </div>
  </CellGroup>

</template>

<script setup>
import { ref } from "vue"
import { CellGroup, Cell, Button, Field, Dialog, Notify, Toast, NoticeBar } from "vant"
import * as api from "../api"
import Qrious from "qrious"

const ownerAddress = ref('')
const showQrTip = ref(false)
const showQr = ref(false)
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
      showQr.value= true
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
  padding: 10px;  
  margin: auto;
  max-width: 600px;
}
.qr-example {
  width: 300px;
  opacity: 0.1;
}
</style>
