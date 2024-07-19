<template>
    <MainNavigation>
        <template #brand-icon>
            <CloudInventoryIcon class="sm:block md:block lg:hidden size-20 p-5" />
        </template>
        <template #account-icon>
            <AccountIcon class="block size-16 p-5" />
        </template>
        <template #aws-logo-icon>
            <AWSLogoIcon class="block size-24 p-5" />
        </template>
        <template #gcp-logo-icon>
            <GCPLogoIcon class="block size-24 p-5" />
        </template>
        <template #azure-logo-icon>
            <AzureLogoIcon class="block size-24 p-5" />
        </template>
    </MainNavigation>
    <HomeContent />
</template>

<script setup>
    import HomeContent from '@/components/HomeContent.vue'
    import MainNavigation from '@/components/MainNavigation.vue'
    import AccountIcon from '@/components/icons/AccountIcon.vue';
    import AWSLogoIcon from '@/components/icons/AWSLogoIcon.vue';
    import AzureLogoIcon from '@/components/icons/AzureLogoIcon.vue';
    import CloudInventoryIcon from '@/components/icons/CloudInventoryIcon.vue';
    import GCPLogoIcon from '@/components/icons/GCPLogoIcon.vue';
</script>