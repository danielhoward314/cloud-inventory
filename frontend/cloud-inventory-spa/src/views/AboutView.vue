<template>
  <div class="about">
    <h1>This is an about page</h1>
    <button @click="getMessage">
      Get Message
    </button>
    <h1>{{ message }}</h1>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const message = ref('')
function getMessage() {
  const fetchOptions = {
    headers: {
      'Content-Type': 'application/json',
    },
    method: 'POST',
    mode: "cors",
    body: JSON.stringify({ type: 'email', email: 'email@here', password: 'antsinmypantsneedtodance' }),
  };
  fetch(`http://localhost:8080/v1/signup`, fetchOptions)
      .then((res) => res.json())
      .then((res) =>{
        console.log(res)
        message.value = res;
      })
}
</script>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
