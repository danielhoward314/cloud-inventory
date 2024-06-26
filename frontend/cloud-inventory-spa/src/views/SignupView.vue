<template>
  <main class="flex h-screen w-screen">
    <SignupForm>
      <template #brand-icon>
        <CloudInventoryIcon class="size-24 p-5" />
      </template>
      <template #password-show-icon>
        <PasswordShowIcon />
      </template>
      <template #password-hide-icon>
        <PasswordHideIcon />
      </template>
    </SignupForm>
  </main>
</template>

<script setup>
  import CloudInventoryIcon from '@/components/icons/CloudInventoryIcon.vue';
  import PasswordHideIcon from '@/components/icons/PasswordHideIcon.vue';
  import PasswordShowIcon from '@/components/icons/PasswordShowIcon.vue';
  import SignupForm from '@/components/SignupForm.vue'
</script>
