<script setup>
import CloudInventoryIcon from '@/components/icons/CloudInventoryIcon.vue';
import PasswordHideIcon from '@/components/icons/PasswordHideIcon.vue';
import PasswordShowIcon from '@/components/icons/PasswordShowIcon.vue';
import LoginForm from '@/components/LoginForm.vue'
</script>

<template>
  <main class="flex h-screen w-screen">
    <LoginForm>
      <template #icon>
        <CloudInventoryIcon class="size-36 p-5" />
      </template>
      <template #password-show-icon>
        <PasswordShowIcon />
      </template>
      <template #password-hide-icon>
        <PasswordHideIcon />
      </template>
    </LoginForm>
  </main>
</template>
