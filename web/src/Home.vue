<template>
  <div>
    <CellGroup>
      <Cell>
        <div class="big-title">Information</div>
      </Cell>
      <Cell title="Device state" icon="info-o" is-link @click="go('/device/state')"></Cell>
      <Cell title="Miner state" icon="info-o" is-link @click="go('/miner/state')" />
    </CellGroup>
  </div>
</template>

<script setup>
import { reactive } from 'vue';
import {
  Tag,
  Col,
  Icon,
  Cell,
  CellGroup,
  Toast,
  Button
} from 'vant';
import { useRouter } from 'vue-router';

const router = useRouter()

function sorry() {
  Toast('暂无后续逻辑~');
}

function go(page) {
  router.push(page)
}

</script>

<style lang="less">
</style>
