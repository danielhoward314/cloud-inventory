<template>
    <TopNav>
        <template #brand-icon>
            <CloudInventoryIcon class="hidden md:block size-20 p-5" />
        </template>
    </TopNav>
</template>

<script setup>
    import TopNav from '@/components/TopNav.vue'
    import CloudInventoryIcon from '@/components/icons/CloudInventoryIcon.vue';
</script>