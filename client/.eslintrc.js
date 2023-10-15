require('@rushstack/eslint-patch/modern-module-resolution');

module.exports = {
  root: true,
  extends: [
    'plugin:vue/vue3-essential',
    '@vue/eslint-config-airbnb',
  ],
  rules: {
    'vuejs-accessibility/form-control-has-label': 'off',
    'vue/max-len': 'off',
  },
};
