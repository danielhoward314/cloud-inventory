import { defineStore } from 'pinia'

export const useAdministratorStore = defineStore('administrator', {
  state: () => ({
    id: '',
    authorizationRole: ''
  }),
  actions: {
    setAdministrator(administrator) {
      if (administrator.id) {
        this.id = administrator.id
      }
      if (administrator.authorizationRole) {
        this.authorizationRole = administrator.authorizationRole
      }
    },
    getAdministrator() {
      return {
        id: this.id,
        authorizationRole: this.authorizationRole
      }
    }
  }
})
