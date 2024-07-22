<template>
  <div class="navbar bg-sky-500 shadow-md">
    <div class="flex-none">
      <CloudInventoryIcon class="size-14 mx-3" />
    </div>
    <div class="flex-none mx-3">
      <RouterLink to="/home" class="btn btn-md btn-square w-20 btn-ghost">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          class="h-10 w-10"
          style="transform: scale(0.6)"
          viewBox="0 0 24 24"
          stroke="white"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
          />
        </svg>
      </RouterLink>
    </div>
    <div class="flex-1">
      <div class="dropdown mx-3">
        <div tabindex="0" role="button" class="btn btn-md btn-square btn-ghost w-20 text-white">
          Resources
        </div>
        <ul tabindex="0" class="dropdown-content menu rounded-box z-[1] p-2 shadow bg-sky-100">
          <li class="my-2">
            <RouterLink
              to="/resources/?provider=aws"
              class="flex items-center justify-start space-x-4"
            >
              <AWSLogoIcon class="block size-10" />
              <span class="text-sky-950">AWS</span>
            </RouterLink>
          </li>
          <li class="my-2">
            <RouterLink
              to="/resources/?provider=gcp"
              class="flex items-center justify-start space-x-4"
            >
              <GCPLogoIcon class="block size-10 transform scale-150" />
              <span class="text-sky-950">GCP</span>
            </RouterLink>
          </li>
          <li class="my-2">
            <RouterLink
              to="/resources/?provider=azure"
              class="flex items-center justify-start space-x-4"
            >
              <AzureLogoIcon class="block size-10" />
              <span class="text-sky-950">Azure</span>
            </RouterLink>
          </li>
        </ul>
      </div>
    </div>
    <div class="flex-none">
      <RouterLink to="/account" class="px-4">
        <svg
          fill="#FFFFFF"
          class="block size-16 p-5"
          style="transform: scale(1.25)"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M9.6,3.32a3.86,3.86,0,1,0,3.86,3.85A3.85,3.85,0,0,0,9.6,3.32M16.35,11a.26.26,0,0,0-.25.21l-.18,1.27a4.63,4.63,0,0,0-.82.45l-1.2-.48a.3.3,0,0,0-.3.13l-1,1.66a.24.24,0,0,0,.06.31l1,.79a3.94,3.94,0,0,0,0,1l-1,.79a.23.23,0,0,0-.06.3l1,1.67c.06.13.19.13.3.13l1.2-.49a3.85,3.85,0,0,0,.82.46l.18,1.27a.24.24,0,0,0,.25.2h1.93a.24.24,0,0,0,.23-.2l.18-1.27a5,5,0,0,0,.81-.46l1.19.49c.12,0,.25,0,.32-.13l1-1.67a.23.23,0,0,0-.06-.3l-1-.79a4,4,0,0,0,0-.49,2.67,2.67,0,0,0,0-.48l1-.79a.25.25,0,0,0,.06-.31l-1-1.66c-.06-.13-.19-.13-.31-.13L19.5,13a4.07,4.07,0,0,0-.82-.45l-.18-1.27a.23.23,0,0,0-.22-.21H16.46M9.71,13C5.45,13,2,14.7,2,16.83v1.92h9.33a6.65,6.65,0,0,1,0-5.69A13.56,13.56,0,0,0,9.71,13m7.6,1.43a1.45,1.45,0,1,1,0,2.89,1.45,1.45,0,0,1,0-2.89Z"
          />
        </svg>
      </RouterLink>
    </div>
  </div>
</template>

<script setup>
import AWSLogoIcon from '@/components/icons/AWSLogoIcon.vue'
import AzureLogoIcon from '@/components/icons/AzureLogoIcon.vue'
import CloudInventoryIcon from '@/components/icons/CloudInventoryIcon.vue'
import GCPLogoIcon from '@/components/icons/GCPLogoIcon.vue'
</script>
