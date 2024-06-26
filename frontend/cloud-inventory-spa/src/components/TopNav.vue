<!-- <template>
    <header class="fixed top-0 w-screen h-12 bg-sky-500 shadow-md">
        <nav class="container mx-auto px-6 py-3">
            <div class="flex justify-start items-center">
                <slot name="brand-icon" />
            </div>
        </nav>
    </header>
</template> -->
<template>
<header class="bg-sky-500 shadow-md w-screen h-auto">
    <nav class="container">
        <div class="flex justify-between items-center">
            <slot name="brand-icon" />
            <div class="flex space-x-4">
                <RouterLink to="/account" class="text-white px-4">Account</RouterLink>
            </div>
        </div>
    </nav>
</header>
</template>

<script setup>
</script>