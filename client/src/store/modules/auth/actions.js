import axios from 'axios';

import environment from '../../../config/environment';
import { SET_USER_DATA, SIGN_IN_ACTION, SIGN_UP_ACTION } from '../../storeconstants';

export default {
  async [SIGN_IN_ACTION](context, payload) {
    try {
      const { status, data } = await axios.post(environment.SIGN_IN_URL, {
        email: payload.email,
        password: payload.password,
      }, {
        validateStatus: () => true,
      });

      if (status !== 200) {
        return data.error;
      }

      const { name, accessToken, refreshToken } = data.data;
      context.commit(SET_USER_DATA, {
        name,
        accessToken,
        refreshToken,
      });

      return null;
    } catch (err) {
      return 'internal server error, try again in a few minutes';
    }
  },
  async [SIGN_UP_ACTION](context, payload) {
    try {
      const { status, data } = await axios.post(environment.SIGN_UP_URL, {
        name: payload.name,
        email: payload.email,
        password: payload.password,
      }, {
        validateStatus: () => true,
      });

      return {
        status,
        data,
      };
    } catch (err) {
      return 'internal server error, try again in a few minutes';
    }
  },
};
