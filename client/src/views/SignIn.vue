<script>
import { mapActions, mapMutations, mapState } from 'vuex';

import { LOADING_SPINNER_MUTATION, SIGN_IN_ACTION } from '../store/storeconstants';
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
      password: '',
    };
  },
  methods: {
    ...mapActions('auth', {
      signin: SIGN_IN_ACTION,
    }),
    ...mapMutations({
      showLoadingMutation: LOADING_SPINNER_MUTATION,
    }),
    async signIn() {
      this.showLoadingMutation(true);
      const error = await this.signin({
        email: this.email,
        password: this.password,
      });
      this.showLoadingMutation(false);

      if (error) {
        this.alertMessage = Array.isArray(error) ? error.join(' | ') : error;
        this.alert.classList.remove('fade');
        return;
      }

      this.email = '';
      this.password = '';

      this.alertMessage = 'you are authenticated :)';
      this.alert.classList.remove('alert-danger');
      this.alert.classList.add('alert-success');
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
          <h3>Sign In</h3>
          <hr>
        </div>

        <form action="">
          <div class="form-group">
            <label for="email">Email</label>
            <input type="text" v-model="email" class="form-control" placeholder="name@email.com">
          </div>
          <div class="form-group">
            <label for="password">Password</label>
            <input type="password" v-model="password" class="form-control">
          </div>

          <div class="my-3 d-grid gap-2">
            <button type="button" v-on:click="signIn" class="btn btn-primary btn-block">Sign In</button>
            <button type="button" class="btn btn-outline-primary mt-3"><a href="/sign_up">Sign Up</a></button>
            <button type="button" class="btn btn-outline-primary"><a href="/forgot_password">Forgot your password?</a></button>
            <button type="button" class="btn btn-outline-primary"><a href="/email_verification_request">Email verification pending?</a></button>
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
