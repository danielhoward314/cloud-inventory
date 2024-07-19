import { defineStore } from 'pinia';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    isAuthenticated: false,
    token: null,
  }),
  actions: {
    setAuthenticated(isAuthenticated, token) {
      this.isAuthenticated = isAuthenticated;
      this.token = token;
    },
    clearAuthentication() {
      this.isAuthenticated = false;
      this.token = null;
    },
  },
});