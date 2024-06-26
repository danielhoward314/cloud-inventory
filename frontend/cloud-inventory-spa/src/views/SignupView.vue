<script setup>
import CloudInventoryIcon from '@/components/icons/CloudInventoryIcon.vue';
import SignupForm from '@/components/SignupForm.vue'
</script>

<template>
  <main class="flex h-screen w-screen">
    <SignupForm>
      <template #icon>
        <CloudInventoryIcon class="size-24 p-5" />
      </template>
    </SignupForm>
  </main>
</template>
