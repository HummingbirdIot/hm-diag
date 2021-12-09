import { createApp } from "vue";
import {
  Tag,
  Col,
  Icon,
  Cell,
  CellGroup,
  Swipe,
  Toast,
  SwipeItem,
  ActionBar,
  ActionBarIcon,
  ActionBarButton,
} from "vant";
import { createRouter, createWebHashHistory } from "vue-router";
import Home from "./Home.vue";
import App from "./App.vue";
import DeviceStateInfo from "./view/DeviceStateInfo.vue";
import MinerStateInfo from "./view/MinerStateInfo.vue";
import Setting from "./view/Setting.vue";
import Control from "./view/Control.vue";
import "./style/common.less";

const routes = [
  { path: "/", component: Home },
  { path: "/setting", component: Setting },
  { path: "/control", component: Control },
  { path: "/device/state", component: DeviceStateInfo },
  { path: "/miner/state", component: MinerStateInfo },
];
const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

const app = createApp(App);

app.use(router).mount("#app");
