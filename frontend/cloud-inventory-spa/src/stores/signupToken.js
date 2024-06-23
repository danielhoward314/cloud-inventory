import { defineStore } from 'pinia';

export const useSignupStore = defineStore('signup', {
  state: () => ({
    hasToken: false,
    token: ''
  }),
  actions: {
    setHasToken(value) {
      this.hasToken = value;
    },
    setToken(value) {
      this.token = value;
    },
    getToken() {
      return this.token
    }
  }
});