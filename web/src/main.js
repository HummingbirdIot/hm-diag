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
import Login from "./Login.vue"
import App from "./App.vue";
import DeviceStateInfo from "./view/DeviceStateInfo.vue";
import MinerStateInfo from "./view/MinerStateInfo.vue";
import Setting from "./view/Setting.vue";
import Control from "./view/Control.vue";
import LogQuery from "./view/LogQuery.vue";
import NetworkTest from "./view/NetworkTest.vue"
import Onboarding from "./view/Onboarding.vue";
import Layout from "./Layout.vue"
import "./style/common.less";
import "./style/index.less";
import { initDateFormat } from "./util/time";
import { AuthToken, setRouter } from './api/auth';
import * as api from './api'

const routes = [
  { path: "/login", component: Login },
  {
    path: "/web",
    component: Layout,
    children: [
      { path: "/", component: Home },
      { path: "/networkTest", component: NetworkTest },
      { path: "/logQuery", component: LogQuery },
      { path: "/setting", component: Setting },
      { path: "/control", component: Control },
      { path: "/onboarding", component: Onboarding },
      { path: "/device/state", component: DeviceStateInfo },
      { path: "/miner/state", component: MinerStateInfo },
    ],
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

router.beforeEach(async (to, from, next) => {
  let configStr = localStorage.config
  let config = {};
  if( configStr ){
    config = JSON.parse(configStr)
  }else{
    config = await api.configGet()
    localStorage.setItem("config",JSON.stringify(config))
  }
  if (config.publicAccess == 1 || config.publicAccess == 0 || AuthToken.get() || window.location.pathname.indexOf("hotspot_tk") == 0) {
    const defaultPage = "/";
    if (to.path === "/login") {
      next({ path: defaultPage });
    } else {
      next();
    }
  } else {
    if (to.path !== "/login") {
      next({ path: "/login" });
    } else {
      next();
    }
  }
});

const app = createApp(App);

app.use(router).use(store).mount("#app");
setRouter(router);
initDateFormat();
