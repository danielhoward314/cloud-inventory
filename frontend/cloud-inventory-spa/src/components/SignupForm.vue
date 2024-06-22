<template>
    <div>
      <h2>Organization Signup</h2>
      <form @submit.prevent="submitForm">
        <div>
          <label for="organizationName">Organization Name:</label>
          <input type="text" id="organizationName" v-model="organizationName" required>
        </div>
        <div>
          <label for="adminEmail">Admin Email:</label>
          <input type="email" id="adminEmail" v-model="adminEmail" required>
        </div>
        <div>
          <label for="adminName">Admin Full Name:</label>
          <input type="text" id="adminName" v-model="adminName" required>
        </div>
        <button type="submit">Sign Up</button>
      </form>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        organizationName: '',
        adminEmail: '',
        adminName: ''
      };
    },
    methods: {
      submitForm() {
        const formData = {
          organization_name: this.organizationName,
          primary_administrator_email: this.adminEmail,
          primary_administrator_name: this.adminName
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
  