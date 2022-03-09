import { createStore } from "vuex";

const store = createStore({
  state() {
    return {
      state: {},
    };
  },
  mutations: {
    state(state, data) {
      state.state = data;
    },
  },
  getters: {
    hasOnboarded: (state) => {
      if (state.state?.miner == null) return null;
      // TODO: more accurate method to judge
      return !!state.state.miner.infoRegion
    },
  },
});

export default store;
