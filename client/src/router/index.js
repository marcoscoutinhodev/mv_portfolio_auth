import { createRouter, createWebHistory } from 'vue-router';

import SignIn from '../views/SignIn.vue';
import SignUp from '../views/SignUp.vue';
import ForgotPassword from '../views/ForgotPassword.vue';
import EmailVerificationRequest from '../views/EmailVerificationRequest.vue';
import ConfirmEmail from '../views/ConfirmEmail.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { name: 'SignIn', path: '/sign_in', component: SignIn },
    { name: 'SignUp', path: '/sign_up', component: SignUp },
    { name: 'ForgotPassword', path: '/forgot_password', component: ForgotPassword },
    { name: 'EmailVerificationRequest', path: '/email_verification_request', component: EmailVerificationRequest },
    { name: 'ConfirmEmail', path: '/confirm_email', component: ConfirmEmail },
  ],
});

export default router;
