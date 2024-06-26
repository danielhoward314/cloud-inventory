<template>
  <div class="md:w-1/2 bg-sky-500 hidden md:flex flex-col">
    <div class="w-full flex">
      <slot name="icon" />
      <div class="flex flex-col justify-center">
        <h1 class="text-lg text-white">Cloud Inventory</h1>
      </div>
    </div>
    <h1 class="text-6xl text-white p-5 h-md:text-lg">A single pane of glass</h1>
    <h1 class="text-6xl text-white p-5 h-md:text-lg">for all your cloud resources</h1>
  </div>
  <div class="w-full md:w-1/2 flex flex-col justify-center items-center h-full">
    <div v-if="hasToken" class="p-5 flex flex-col justify-center items-center space-y-4">
      <h2 class="text-2xl text-sky-950 font-semibold tracking-tight">Verify Your Email Address</h2>
      <p class="text-lg text-sky-950 font-light tracking-tight">Enter verification code below</p>
      <form @submit.prevent="handleSubmit" ref="form" class="space-y-4">
        <div>
          <label for="verificationCode" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 sr-only">Verification Code</label>
          <input type="text" id="verificationCode" v-model="verificationCode" @input="checkCodeLength" required class="flex h-9 sm:w-48 md:w-96 rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50">
        </div>
      </form>
    </div>
    <div v-else class="p-5 flex flex-col justify-center items-center space-y-4">
      <h2 class="text-2xl text-sky-950 font-semibold tracking-tight">Create an account</h2>
      <form @submit.prevent="submitForm" class="space-y-4">
        <div>
          <label for="organizationName" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 sr-only">Organization Name</label>
          <CustomInput
            placeholder="Organization Name"
            id="organizationName"
            v-model="organizationName"
            customClass="flex h-9 sm:w-48 md:w-96 rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
            required
          />
        </div>
        <div>
          <label for="adminName" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 sr-only">Name</label>
          <CustomInput
            placeholder="Name"
            id="adminName"
            v-model="adminName"
            customClass="flex h-9 sm:w-48 md:w-96 rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
            required
          />
        </div>
        <div>
          <label for="adminEmail" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 sr-only">Email</label>
          <CustomInput
            type="email"
            placeholder="name@email.com"
            id="adminEmail"
            v-model="adminEmail"
            customClass="flex h-9 sm:w-48 md:w-96 rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
            required
          />
        </div>
        <div>
          <label for="password" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 sr-only">Password</label>
          <CustomInput
            type="text"
            placeholder="Password"
            id="password"
            v-model="password"
            customClass="flex h-9 sm:w-48 md:w-96 rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
            required
          />
        </div>
        <button type="submit" class="inline-flex sm:w-48 md:w-96 items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 bg-sky-950 text-white shadow hover:bg-sky-900/90 h-9 px-4 py-2">Sign Up</button>
      </form>
    </div>
  </div>
</template>
  
<script>
import { ref, watch, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useSignupStore } from '@/stores/signupToken';
import router from '@/router';
import constants from '@/consts/consts';
import CustomInput from '@/components/reusable/CustomInput.vue';

export default {
  components: {
    CustomInput,
  },
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
        if (data.jwt) {
          localStorage.setItem(constants.localStorageKeys.sessionJWT, data.jwt)
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
  