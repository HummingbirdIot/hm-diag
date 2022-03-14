import { createStore } from "vuex";

const store = createStore({
  state() {
    return {
      state: {}, // hotspot state
      safeConf: {},
      envState: {
        isInLocal: false,
        isViaPrivate: false,
      },
    };
  },
  mutations: {
    state(state, data) {
      state.state = data;
    },
    isInLocal(state, v) {
      state.envState.isInLocal = v;
    },
    isViaPrivate(state, v) {
      state.envState.isViaPrivate = v;
    },
    safeConf(state, data) {
      state.safeConf = data;
    },
  },
  getters: {
    hasOnboarded: (state) => {
      if (state.state?.miner == null) return null;
      // TODO: more accurate method to judge
      return !!state.state.miner.infoRegion;
    },
    canAccessImportant: (state, getters) => {
      // server get remote ip is private ip
      if (state.envState.isViaPrivate) return true;
      // server config: can public access 
      if (state.safeConf.publicAccess != null) {
        return state.safeConf.publicAccess == 1;
      }
      // browser location address is private ip
      return store.state.isInLocal;
    },
  },
});

export default store;
