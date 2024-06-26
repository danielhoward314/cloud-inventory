<template>
    <div class="h-screen w-screen flex flex-col justify-start items-center space-y-4">
      <slot name="icon" />
      <h1 class="text-lg text-sky-950 font-light">Log in to Cloud Inventory</h1>
      <div class="flex flex-col lg:w-200 h-auto bg-sky-500 border-sky-950 rounded-md p-4 shadow-sm">
        <form @submit.prevent="submitForm" class="space-y-4">
          <div>
            <label for="email" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 sr-only">Email</label>
            <CustomInput
              type="email"
              placeholder="name@yourorganization.com"
              id="email"
              v-model="form.email"
              customClass="flex h-9 sm:w-48 md:w-72 lg:w-96 focus:bg-slate-100 rounded-md border border-input bg-white px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
              required
            />
          </div>
          <div>
            <label for="password" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 sr-only">Password</label>
            <CustomInput
              type="password"
              placeholder="Password"
              id="password"
              v-model="form.password"
              customClass="flex h-9 sm:w-48 md:w-72 lg:w-96 focus:bg-slate-100 rounded-md border border-input bg-white px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
              required
            />
          </div>
          <button type="submit" class=" flex sm:w-48 md:w-72 lg:w-96 items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 bg-sky-950 text-white shadow hover:bg-sky-900/90 h-9 px-4 py-2">Log In</button>
        </form>
      </div>
      <div class="flex flex-col justify-center items-center lg:w-200 h-auto border-sky-500 border-2 rounded-md p-4 shadow-md">
        <div>
          <p class="h-9 sm:w-48 md:w-72 lg:w-96 md:text-md lg:text-lg text-sky-950 text-center font-light">New to Cloud Inventory?</p>
        </div>
        <div>
          <RouterLink to="/signup" class="flex justify-center items-center h-9 sm:w-48 md:w-72 lg:w-96 md:text-md lg:text-lg text-sky-500 text-center font-light">Create an account</RouterLink>
        </div>
      </div>
    </div>
  </template>
  
  <script>
import constants from '@/consts/consts';
import router from '@/router';
import CustomInput from '@/components/reusable/CustomInput.vue';

  export default {
    components: {
      CustomInput,
    },
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
  