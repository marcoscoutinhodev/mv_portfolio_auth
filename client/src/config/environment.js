export default {
  SIGN_IN_URL: import.meta.env.SIGN_IN_URL || 'http://localhost:8080/auth/signin',
  SIGN_UP_URL: import.meta.env.SIGN_UP_URL || 'http://localhost:8080/auth/signup',
  FORGOT_PASSWORD_URL: import.meta.env.FORGOT_PASSWORD_URL || 'http://localhost:8080/auth/forgot-password',
  EMAIL_CONFIRMATION_REQUEST_URL: import.meta.env.EMAIL_CONFIRMATION_REQUEST_URL || 'http://localhost:8080/auth/email-confirmation-request',
  CONFIRM_EMAIL_REQUEST_URL: import.meta.env.CONFIRM_EMAIL_REQUEST_URL || 'http://localhost:8080/auth/confirm-email',
};
