<template>
  <div>Account</div>
  <RouterLink to="/home" class="px-4">
    <h1>Home</h1>
  </RouterLink>
  <AccountContent />
</template>

<script setup>
import AccountContent from '@/components/AccountContent.vue'
</script>
