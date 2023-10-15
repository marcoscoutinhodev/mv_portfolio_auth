<script>
import { mapActions, mapMutations, mapState } from 'vuex';

import { LOADING_SPINNER_MUTATION, SIGN_UP_ACTION } from '../store/storeconstants';
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
      name: '',
      email: '',
      password: '',
    };
  },
  methods: {
    ...mapActions('auth', {
      signup: SIGN_UP_ACTION,
    }),
    ...mapMutations({
      showLoadingMutation: LOADING_SPINNER_MUTATION,
    }),
    async signUp() {
      this.alert.classList.remove('alert-successful');
      this.alert.classList.add('alert-danger');
      this.alert.classList.add('fade');

      this.showLoadingMutation(true);
      const { status, data } = await this.signup({
        name: this.name,
        email: this.email,
        password: this.password,
      });
      this.showLoadingMutation(false);

      if (status !== 201) {
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
          <h3>Sign Up</h3>
          <hr>
        </div>

        <form action="">
          <div class="form-group">
            <label for="name">Name</label>
            <input type="text" class="form-control" placeholder="Full Name" v-model="name">
          </div>
          <div class="form-group">
            <label for="email">Email</label>
            <input type="text" class="form-control" placeholder="name@email.com" v-model="email">
            <div id="emailHelper" class="form-text">We'll never share your email with anyone else</div>
          </div>
          <div class="form-group">
            <label for="password">Password</label>
            <input type="password" class="form-control" v-model="password">
            <small id="passwordHelpBlock" class="form-text text-muted">
              Your password must be at least 7 characters long, must contain special characters "!@#$%^&*(),.?":{}|&lt;>", numbers, lowercase and uppercase letters only.
            </small>
          </div>

          <div class="my-3 d-grid gap-2">
            <button type="button" v-on:click="signUp" class="btn btn-primary btn-block">Sign Up</button>
            <button type="button" class="btn btn-outline-primary mt-3"><a href="/sign_in">Sign In</a></button>
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
