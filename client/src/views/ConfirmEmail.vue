<script>
import { mapActions, mapMutations } from 'vuex';

import { LOADING_SPINNER_MUTATION, CONFIRM_EMAIL_REQUEST_ACTION } from '../store/storeconstants';
import TheLoader from '../components/TheLoader.vue';

export default {
  components: {
    TheLoader,
  },
  data() {
    return {
      alert: null,
      alertMessage: 'invalid credentials',
    };
  },
  methods: {
    ...mapActions('auth', {
      confirmEmail: CONFIRM_EMAIL_REQUEST_ACTION,
    }),
    ...mapMutations({
      showLoadingMutation: LOADING_SPINNER_MUTATION,
    }),
  },
  mounted() {
    this.alert = document.querySelector('#alert');
  },
  async created() {
    this.showLoadingMutation(true);
    const token = this.$route.query.t;

    const { status, data } = await this.confirmEmail({
      token,
    });
    this.showLoadingMutation(false);

    if (status !== 200) {
      this.alertMessage = Array.isArray(data.error) ? data.error.join(' | ') : data.error;
      this.alert.classList.remove('fade');
      return;
    }

    this.alertMessage = 'email confirmed';
    this.alert.classList.remove('alert-danger');
    this.alert.classList.add('alert-success');
    this.alert.classList.remove('fade');
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
          <h3>Confirming email</h3>
          <hr>
        </div>

        <div id="alert" class="alert alert-danger fade" role="alert">
          {{ alertMessage }}
        </div>
      </div>
    </div>
  </div>

</template>
