<template>
  <router-view/>
  <Tabbar route>
    <TabbarItem replace to="/" icon="home-o" active>Home</TabbarItem>
    <TabbarItem replace to="/setting" icon="setting-o">Setting</TabbarItem>
    <TabbarItem replace to="/control" icon="diamond-o">Control</TabbarItem>
    <TabbarItem v-if="showNeighbor" replace to="/neighbor" icon="search">Neighbors</TabbarItem>
  </Tabbar>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useStore } from 'vuex'
import { Tabbar, TabbarItem } from 'vant';
import * as api from './api'
import isPrivateIp from 'private-ip'

const store = useStore()
const showNeighbor = ref(true)

onMounted(()=>{
  api.stateGet().then(s=>{
    store.commit('state', s)
  }).catch(e=>{
    console.error('init hotspot state data error: ', e)
  })
  localEnvJudge()
  isViaPrivate()
})

function localEnvJudge() {
  const h = window.location.hostname
  if (h === 'localhost') return true
  const r = isPrivateIp(h)
  store.commit('isInLocal', r)
}

function isViaPrivate() {
  api.isViaPrivate()
    .then(r=>{
      store.commit('isViaPrivate', r)
    })
}

function isShowNeighbors() {
  if (window.location.pathname.indexOf('/proxy/') === 0) {
    showNeighbor.value = false
  }
}
isShowNeighbors()


</script>

