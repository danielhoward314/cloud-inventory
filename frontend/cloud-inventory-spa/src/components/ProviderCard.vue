<template>
    <div>
        <div v-for="provider in providers" :key="provider.id" class="card bg-base-100 w-96 shadow-xl">
        <figure class="px-10 pt-10">
            <img :src="provider.image" :alt="provider.title" class="rounded-xl" />
        </figure>
        <div class="card-body items-center text-center">
            <h2 class="card-title">{{ provider.title }}</h2>
            <p>{{ provider.description }}</p>
            <div class="card-actions">
            <button class="btn btn-primary">{{ provider.buttonText }}</button>
            </div>
        </div>
        </div>
    </div>
</template>
  
<script>
  export default {
    name: 'ProviderCard',
    props: {
      providers: {
        type: Array,
        required: true,
        default: () => []
      }
    }
  }
</script>
