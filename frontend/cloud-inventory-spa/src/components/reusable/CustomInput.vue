<template>
  <input
    :type="type"
    :placeholder="placeholder"
    :id="id"
    v-model="localValue"
    :required="required"
    :class="computedClass"
    @input="$emit('update:modelValue', localValue)"
  />
</template>

<script>
export default {
  name: 'CustomInput',
  props: {
    type: {
      type: String,
      default: 'text'
    },
    placeholder: {
      type: String,
      default: ''
    },
    id: {
      type: String,
      required: true
    },
    modelValue: {
      type: [String, Number],
      default: ''
    },
    required: {
      type: Boolean,
      default: false
    },
    customClass: {
      type: String,
      default: ''
    }
  },
  computed: {
    localValue: {
      get() {
        return this.modelValue
      },
      set(value) {
        this.$emit('update:modelValue', value)
      }
    },
    computedClass() {
      return this.customClass
    }
  }
}
</script>

<style scoped></style>
