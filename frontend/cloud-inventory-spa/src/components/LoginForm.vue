<template>
    <div class="login-form">
      <h2>Cloud Inventory</h2>
      <form @submit.prevent="submitForm">
        <div class="form-group">
          <label for="email">Email</label>
          <input
            type="text"
            id="email"
            v-model="form.email"
            required
          />
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <input
            type="password"
            id="password"
            v-model="form.password"
            required
          />
        </div>
        <button type="submit">Login</button>
      </form>
    </div>
  </template>
  
  <script>
import constants from '@/consts/consts';
import router from '@/router';

  export default {
    data() {
      return {
        form: {
          email: '',
          password: '',
        },
      };
    },
    methods: {
      async submitForm() {
        const formData = {
        email: this.form.email,
        password: this.form.password,
      };
      const fetchOptions = {
          headers: {
          'Content-Type': 'application/json',
          },
          method: 'POST',
          mode: "cors",
          body: JSON.stringify(formData),
      };
      fetch('http://localhost:8080/v1/login', fetchOptions)
      .then(response => response.json())
      .then(data => {
        if (data && data.jwt) {
          localStorage.setItem(constants.localStorageKeys.sessionJWT, data.jwt);
          router.push('/home');
        } else {
          console.error('No login response data received.');
        }
      })
      .catch(error => {
        console.error('Error submitting form:', error);
      });
      },
    },
  };
  </script>
  
  <style scoped>
  </style>
  