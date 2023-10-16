<script>
import { mapActions, mapMutations, mapState } from 'vuex';

import { LOADING_SPINNER_MUTATION, UPDATE_PASSWORD_ACTION } from '../store/storeconstants';
import TheLoader from '../components/TheLoader.vue';

export default {
  components: {
    TheLoader,
  },
  mounted() {
    this.alert = document.querySelector('#alert');
  },
  data() {
    return {
      alert: null,
      alertMessage: 'invalid credentials',
      password: '',
    };
  },
  methods: {
    ...mapActions('auth', {
      updatepassword: UPDATE_PASSWORD_ACTION,
    }),
    ...mapMutations({
      showLoadingMutation: LOADING_SPINNER_MUTATION,
    }),
    async updatePassword() {
      this.showLoadingMutation(true);
      const token = this.$route.query.t;
      const { status, data } = await this.updatepassword({
        password: this.password,
        token,
      });
      this.showLoadingMutation(false);

      if (status !== 200) {
        this.alertMessage = Array.isArray(data.error) ? data.error.join(' | ') : data.error;
        this.alert.classList.remove('fade');
        return;
      }

      this.password = '';

      this.alertMessage = 'Password updated successfully';
      this.alert.classList.remove('alert-danger');
      this.alert.classList.add('alert-success');
      this.alert.classList.remove('fade');

      setTimeout(() => {
        this.$router.push('/sign_in');
      }, 1500);
    },
  },
  computed: {
    ...mapState({
      showLoading: (state) => state.showLoading,
    }),
  },
};

</script>

<template>
  <TheLoader v-if="showLoading" />
  <div class="row">
    <div class="col-md-4 offset-md-4">
      <div>

        <div class="container">
          <div class="col-md-4 px-0 rounded mx-auto d-block mt-5">
            <img src="../assets/logo/png/logo-no-background.png" alt="logo" class="img-fluid">
          </div>
        </div>

        <div class="mt-3 text-center">
          <h3>Update your password</h3>
          <hr>
        </div>

        <form action="">
          <div class="form-group">
            <label for="password">Passord</label>
            <input type="password" v-model="password" class="form-control">
          </div>

          <div class="my-3 d-grid gap-2">
            <button type="button" v-on:click="updatePassword" class="btn btn-primary btn-block">Update Password</button>
          </div>

          <div class="form-group">
            <div id="alert" class="alert alert-danger fade" role="alert">
              {{ alertMessage }}
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>

</template>
