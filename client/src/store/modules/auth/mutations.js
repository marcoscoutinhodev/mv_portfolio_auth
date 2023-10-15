import { SET_USER_DATA } from '../../storeconstants';

export default {
  [SET_USER_DATA](state, payload) {
    state.name = payload.name;
    state.accessToken = payload.accessToken;
    state.refreshToken = payload.refreshToken;
  },
};
