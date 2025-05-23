<template>
  <div
    ref="root"
    class="relative overflow-hidden"
    :class="customClass"
  >
    <div class="h-full w-full rounded-[inherit] overflow-auto pr-2" ref="viewport">
      <slot />
    </div>
    <!-- <div class="absolute right-0 top-0 h-full w-2.5 border-l border-l-transparent p-[1px]">
      <div class="relative h-full w-full rounded-full bg-border" />
    </div> -->
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

defineProps({
  customClass: {
    type: String,
    default: ''
  }
})

const root = ref(null)
const viewport = ref(null)
</script>

<style scoped>
/* You can optionally enhance with scrollbar styling here */
.bg-border {
  background-color: rgba(0, 0, 0, 0.2); /* Adjust for theme */
}
</style>
