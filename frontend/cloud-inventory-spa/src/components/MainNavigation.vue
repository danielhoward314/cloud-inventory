<template>
    <div class="navbar bg-sky-500 shadow-md">
        <div class="flex-none">
            <RouterLink to="/home" class="px-4">
                <slot name="brand-icon" />
                <h1 class="hidden lg:block text-white text-center px-4">Cloud Inventory</h1>
            </RouterLink>
        </div>
        <div class="flex-1 drawer">
            <input id="main-drawer" type="checkbox" class="drawer-toggle" />
            <div class="drawer-content">
                <label for="main-drawer" class="btn btn-square btn-ghost drawer-button">
                    <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                class="inline-block h-5 w-5 stroke-current">
                <path
                    stroke="white"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M4 6h16M4 12h16M4 18h16"
                />
                    </svg>
                </label>
            </div>
            <div class="drawer-side">
                <label for="main-drawer" aria-label="close sidebar" class="drawer-overlay"></label>
                <ul class="menu bg-base-200 text-base-content min-h-full w-80 p-4">
                    <li>
                        <RouterLink to="/aws" class="px-4">
                            <slot name="aws-logo-icon" />
                        </RouterLink>
                    </li>
                    <li>
                        <RouterLink to="/gcp" class="px-4">
                            <slot name="gcp-logo-icon" />
                        </RouterLink>
                    </li>
                    <li>
                        <RouterLink to="/azure" class="px-4">
                            <slot name="azure-logo-icon" />
                        </RouterLink>
                    </li>
                </ul>
            </div>
        </div>
        <div class="flex-none">
            <RouterLink to="/account" class="px-4">
                <slot name="account-icon" />
            </RouterLink>
        </div>
    </div>
</template>

<script setup>
</script>