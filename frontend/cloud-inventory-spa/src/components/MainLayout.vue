<template>
  <div>
    <MainNavigation />
    <main>
      <slot></slot>
    </main>
  </div>
</template>

<script setup>
import MainNavigation from '@/components/MainNavigation.vue'
</script>
