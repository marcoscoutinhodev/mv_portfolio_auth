import { createStore } from 'vuex';

import auth from './modules/auth';
import { LOADING_SPINNER_MUTATION } from './storeconstants';

const store = createStore({
  modules: {
    auth,
  },
  state() {
    return {
      showLoading: false,
    };
  },
  mutations: {
    [LOADING_SPINNER_MUTATION](state, payload) {
      state.showLoading = payload;
    },
  },
});

export default store;
