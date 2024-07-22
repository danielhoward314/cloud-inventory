<template>
  <section class="bg-sky-100 w-screen h-screen flex flex-col items-center p-4 space-y-8">
    <BaseCard class="w-1/2" bg="bg-sky-200">
      <div class="flex flex-col justify-center items-center">
        <h2 class="text-2xl font-bold text-sky-950">Add Provider Connection</h2>
        <p class="my-4 pb-4 text-sky-950 font-light">
          Grant access to your cloud provider accounts.
        </p>
      </div>
      <div class="flex justify-around">
        <div class="flex btn w-1/3 btn-ghost items-center justify-center">
          <AWSLogoIcon class="size-10 py-1 transform scale-110" />
        </div>
        <div class="flex w-1/3 btn btn-ghost items-center justify-center">
          <GCPLogoIcon class="size-10 transform scale-150" />
        </div>
        <div class="flex w-1/3 btn btn-ghost items-center justify-center">
          <AzureLogoIcon class="size-10" />
        </div>
      </div>
    </BaseCard>
    <div v-if="error" class="toast">
  <div class="alert alert-error text-sky-100 text-wrap">
    <span>Error loading provider connections</span>
    <button @click="closeErrToast" class="btn btn-ghost">
    <svg  class="w-3 h-3 hover:" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"/>
        </svg>
      </button>
  </div>
</div>
    <BaseCard class="w-1/2" bg="bg-sky-200">
      <DataTable :value="rows" :loading="loading" showGridlines class="text-sky-950 bg-sky-100">
        <Column
          v-for="(col, idx) of columns"
          :key="idx"
          :field="col.field"
          :header="col.header"
        ></Column>
      </DataTable>
    </BaseCard>
  </section>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { useOrganizationStore } from '@/stores/organization'
import xhrClient from '@/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import BaseCard from '@/components/reusable/BaseCard.vue'
import AWSLogoIcon from '@/components/icons/AWSLogoIcon.vue'
import AzureLogoIcon from '@/components/icons/AzureLogoIcon.vue'
import GCPLogoIcon from '@/components/icons/GCPLogoIcon.vue'
import constants from '@/consts/consts'

const loading = ref(true)
const error = ref(false)
const rows = ref()

const columns = ref([
  { header: 'Name', field: 'display_name' },
  { header: 'Account Identifier', field: 'externalIdentifier' },
  { header: 'Connection', field: 'connection' }
])

onMounted(async () => {
  try {
    const orgStore = useOrganizationStore()
    const org = orgStore.getOrganization()
    const res = await xhrClient.getProviders(org.id)
    if (res && res.data && res.data.providers) {
      rows.value = res.data.providers.map((provider) => {
        const row = {
          display_name: provider.name,
          externalIdentifier: provider.externalIdentifier
        }
        if (provider.providerName === constants.providers.awsProviderName) {
          row.connection = provider.awsMetadata.roleArn
        } else if (provider.providerName === constants.providers.gcpProviderName) {
          row.connection = provider.gcpMetadata.serviceAccountId
        } else {
          row.connection = provider.azureMetadata.serviceAccountId
        }
        return row
      })
    }
    loading.value = false
  } catch (e) {
    console.error('Error fetching provider details:', e)
    loading.value = false
    error.value = true
  }
})

const closeErrToast = () => {
  error.value = false
}

onUnmounted(() => {
  loading.value = true
  error.value = false
  rows.value = null
})
</script>
