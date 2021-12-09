<template>
  <router-view/>
  <Tabbar route>
    <TabbarItem replace to="/" icon="home-o" active>Home</TabbarItem>
    <TabbarItem replace to="/setting" icon="setting-o">Setting</TabbarItem>
    <TabbarItem replace to="/control" icon="diamond-o">Control</TabbarItem>
  </Tabbar>
</template>

<script setup>
import { Tabbar, TabbarItem } from 'vant';

</script>

