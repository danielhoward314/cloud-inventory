<template>
  <section class="w-screen h-screen bg-sky-100">
    <ProviderCard :providers="providers" />
  </section>
</template>

<script>
import xhrClient from '@/api'
import ProviderCard from '@/components/ProviderCard.vue'

export default {
  components: {
    ProviderCard
  },
  data() {
    return {
      providers: []
    }
  },
  mounted() {
    this.getProviders()
  },
  methods: {
    async getProviders() {
      try {
        const res = await xhrClient.getProviders({ test: 'test data here' })
        console.log(res.data)
        const resTwo = await xhrClient.getProvidersExtraConfig(
          { test: 'test data here' },
          { params: { id: 'madeup' } }
        )
        console.log(resTwo.data)
        // this.providers = data.map((provider) => ({
        //   id: provider.id,
        //   image: provider.image_url,
        //   title: provider.name,
        //   description: provider.description,
        //   buttonText: provider.button_text
        // }))
      } catch (error) {
        console.error('Error fetching providers:', error)
      }
    }
  }
}
</script>
