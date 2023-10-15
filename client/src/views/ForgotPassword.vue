<script>
import { mapActions, mapMutations, mapState } from 'vuex';

import { LOADING_SPINNER_MUTATION, FORGOT_PASSWORD_ACTION } from '../store/storeconstants';
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
      email: '',
    };
  },
  methods: {
    ...mapActions('auth', {
      forgotpassword: FORGOT_PASSWORD_ACTION,
    }),
    ...mapMutations({
      showLoadingMutation: LOADING_SPINNER_MUTATION,
    }),
    async forgotPassword() {
      this.showLoadingMutation(true);
      const { status, data } = await this.forgotpassword({
        email: this.email,
      });
      this.showLoadingMutation(false);

      if (status !== 200) {
        this.alertMessage = Array.isArray(data.error) ? data.error.join(' | ') : data.error;
        this.alert.classList.remove('fade');
        return;
      }

      this.name = '';
      this.email = '';
      this.password = '';

      this.alertMessage = data.data;
      this.alert.classList.remove('alert-danger');
      this.alert.classList.add('alert-successful');
      this.alert.classList.remove('fade');
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
          <h3>Forgot your password?</h3>
          <hr>
        </div>

        <form action="">
          <div class="form-group">
            <label for="email">Email</label>
            <input type="text" v-model="email" class="form-control" placeholder="name@email.com">
          </div>

          <div class="my-3 d-grid gap-2">
            <button type="button" v-on:click="forgotPassword" class="btn btn-primary btn-block">Recover Password</button>
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
