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
import store from "./store"
import Home from "./Home.vue";
import App from "./App.vue";
import DeviceStateInfo from "./view/DeviceStateInfo.vue";
import MinerStateInfo from "./view/MinerStateInfo.vue";
import Setting from "./view/Setting.vue";
import Control from "./view/Control.vue";
import LogQuery from "./view/LogQuery.vue";
import Neighbor from "./view/Neighbor.vue";
import Onboarding from "./view/Onboarding.vue";
import "./style/common.less";
import "./style/index.less";
import { initDateFormat } from "./util/time";

const routes = [
  { path: "/", component: Home },
  { path: "/logQuery", component: LogQuery },
  { path: "/setting", component: Setting },
  { path: "/control", component: Control },
  { path: "/onboarding", component: Onboarding },
  { path: "/neighbor", component: Neighbor},
  { path: "/device/state", component: DeviceStateInfo },
  { path: "/miner/state", component: MinerStateInfo },
];
const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

const app = createApp(App);

app.use(router).use(store).mount("#app");

initDateFormat();
