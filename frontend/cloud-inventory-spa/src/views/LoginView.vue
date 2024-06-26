<script setup>
import CloudInventoryIcon from '@/components/icons/CloudInventoryIcon.vue';
import LoginForm from '@/components/LoginForm.vue'
</script>

<template>
  <main class="flex h-screen w-screen">
    <LoginForm>
      <template #icon>
        <CloudInventoryIcon class="size-36 p-5" />
      </template>
    </LoginForm>
  </main>
</template>
