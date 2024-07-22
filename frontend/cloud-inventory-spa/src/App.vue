<template>
  <component :is="layout">
    <RouterView />
  </component>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import MainLayout from './components/MainLayout.vue'
import BlankLayout from './components/BlankLayout.vue'
const route = useRoute()
const layout = computed(() => (route.meta.layout === 'blank' ? BlankLayout : MainLayout))
</script>
