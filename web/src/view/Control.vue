<template>
  <Cell>
    <div class="big-title">Control</div>
  </Cell>

  <CellGroup>
    <Cell title="Reboot device">
      <Button size="small" type="danger">reboot</Button>
    </Cell>
    <Cell title="Resync miner">
      <Button size="small" type="primary">Resync</Button>
    </Cell>
  </CellGroup>

  <CellGroup>
    <Cell title="Generate snapshot">
      <Button size="small" type="primary">Download</Button>
    </Cell>
    <Cell title="Upload snapshot">
      <Button size="small" type="primary">Upload</Button>
    </Cell>
  </CellGroup>
</template>

<script setup>
import { CellGroup, Cell, Field, Button, NoticeBar } from 'vant'

</script>

