<template>
    <div v-if="hasToken">
      <h2>Verify Your Email Address</h2>
      <h3>Enter verification code below</h3>
      <form @submit.prevent="handleSubmit" ref="form">
        <div>
          <label for="verificationCode">Verification Code:</label>
          <input type="text" id="verificationCode" v-model="verificationCode" @input="checkCodeLength" required>
        </div>
      </form>
    </div>
    <div v-else>
      <h2>Organization Signup</h2>
      <form @submit.prevent="submitForm">
        <div>
          <label for="organizationName">Organization Name:</label>
          <input type="text" id="organizationName" v-model="organizationName" required>
        </div>
        <div>
          <label for="adminName">Admin Full Name:</label>
          <input type="text" id="adminName" v-model="adminName" required>
        </div>
        <div>
          <label for="adminEmail">Admin Email:</label>
          <input type="email" id="adminEmail" v-model="adminEmail" required>
        </div>
        <div>
          <label for="password">Password:</label>
          <input type="text" id="password" v-model="password" required>
        </div>
        <button type="submit">Sign Up</button>
      </form>
    </div>
  </template>
  
<script>
import { ref, watch, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useSignupStore } from '../stores/signupToken';
import router from '@/router';

export default {
  setup() {
    const route = useRoute();
    const signupStore = useSignupStore();
    const token = ref('');
    const verificationCode = ref('');

    const checkToken = () => {
      if (route.query.token) {
        signupStore.setHasToken(true);
        signupStore.setToken(route.query.token);
      } else {
        signupStore.setHasToken(false);
        signupStore.setToken('');
      }
    };

    const checkCodeLength = () => {
      if (verificationCode.value.length === 6) {
        submitCodeVerificationForm();
      }
    };

    const submitCodeVerificationForm = () => {
      const formData = {
        token: signupStore.getToken(),
        verification_code: verificationCode.value,
      };
      const fetchOptions = {
          headers: {
          'Content-Type': 'application/json',
          },
          method: 'POST',
          mode: "cors",
          body: JSON.stringify(formData),
      };
      fetch('http://localhost:8080/v1/verify', fetchOptions)
      .then(response => response.json())
      .then(data => {
        if (data.success) {
          router.push('/home')
        }
      })
      .catch(error => {
        console.error('Error submitting form:', error);
      });
    };

    watch(route, () => {
      checkToken();
    }, { immediate: true });

    onMounted(() => {
      checkToken();
    });

    return {
      checkCodeLength,
      hasToken: signupStore.hasToken,
      submitCodeVerificationForm,
      token,
      verificationCode,
    };
  },
  data() {
    return {
      organizationName: '',
      adminEmail: '',
      adminName: '',
      password: '',
    };
  },
  methods: {
    submitForm() {
      const formData = {
        organization_name: this.organizationName,
        primary_administrator_email: this.adminEmail,
        primary_administrator_name: this.adminName,
        primary_administrator_cleartext_password: this.password,
      };
      const fetchOptions = {
          headers: {
          'Content-Type': 'application/json',
          },
          method: 'POST',
          mode: "cors",
          body: JSON.stringify(formData),
      };
      fetch('http://localhost:8080/v1/signup', fetchOptions)
      .then(response => response.json())
      .then(data => {
        if (data.token) {
          window.location.href = `/signup?token=${data.token}`;
        } else {
          console.error('No token received in the response.');
        }
      })
      .catch(error => {
        console.error('Error submitting form:', error);
      });
    }
  }
};
</script>
  
  <style scoped>
  /* Add your component styles here */
  </style>
  