<template>
  <div>
    <div v-if="id">ID: {{ id }}</div>
    <div v-if="name">Name: {{ name }}</div>
    <div v-if="billingPlan">Billing Plan: {{ billingPlan }}</div>
    <div v-else>No organization data available</div>
  </div>
</template>

<script>
import { mapState } from 'pinia'
import xhrClient from '@/api'
import { useOrganizationStore } from '@/stores/organization'

export default {
  computed: {
    ...mapState(useOrganizationStore, ['id', 'name', 'billingPlan'])
  },
  mounted() {
    this.getOrganization()
  },
  methods: {
    async getOrganization() {
      try {
        const orgStore = useOrganizationStore()
        const org = orgStore.getOrganization()
        const res = await xhrClient.getOrganization(org.id)
        orgStore.setOrganization({
          id: res.data.id,
          name: res.data.organizationName,
          billingPlan: res.data.billingPlan
        })
      } catch (error) {
        console.error('Error fetching account details:', error)
      }
    }
  }
}
</script>
