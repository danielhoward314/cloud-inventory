import { defineStore } from 'pinia'

export const useOrganizationStore = defineStore('organization', {
  state: () => ({
    id: '',
    name: '',
    billingPlan: ''
  }),
  actions: {
    setOrganization(organization) {
      if (organization.id) {
        this.id = organization.id
      }
      if (organization.name) {
        this.name = organization.name
      }
      if (organization.billingPlan) {
        this.billingPlan = organization.billingPlan
      }
    },
    getOrganization() {
      return {
        id: this.id,
        name: this.name,
        billingPlan: this.billingPlan
      }
    }
  }
})
