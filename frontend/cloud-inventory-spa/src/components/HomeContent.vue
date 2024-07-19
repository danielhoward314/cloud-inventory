<template>
    <section class="w-screen h-screen bg-sky-100">
      <ProviderCard :providers="providers" />
    </section>
</template>

<script>
import ProviderCard from '@/components/ProviderCard.vue';

export default {
  components: {
    ProviderCard
  },
  data() {
    return {
      providers: []
    };
  },
  mounted() {
    this.fetchProviders();
  },
  methods: {
    async fetchProviders() {
      try {
        const response = await fetch('https://api.example.com/providers');
        const data = await response.json();
        this.providers = data.map(provider => ({
          id: provider.id,
          image: provider.image_url,
          title: provider.name,
          description: provider.description,
          buttonText: provider.button_text
        }));
      } catch (error) {
        console.error('Error fetching providers:', error);
      }
    }
  }
}
</script>