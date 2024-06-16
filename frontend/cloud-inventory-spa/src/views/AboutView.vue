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
  fetch("http://localhost:8080/v0/hello")
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
