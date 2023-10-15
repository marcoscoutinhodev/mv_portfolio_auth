import { createRouter, createWebHistory } from 'vue-router';

import SignIn from '../views/SignIn.vue';
import SignUp from '../views/SignUp.vue';
import ForgotPassword from '../views/ForgotPassword.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { name: 'SignIn', path: '/sign_in', component: SignIn },
    { name: 'SignUp', path: '/sign_up', component: SignUp },
    { name: 'ForgotPassword', path: '/forgot_password', component: ForgotPassword },
  ],
});

export default router;
