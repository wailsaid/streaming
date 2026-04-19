<script setup lang="ts">
import Plyr from 'plyr'
import 'plyr/dist/plyr.css'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = defineProps<{
  src?: string
  poster?: string
}>()

const videoRef = ref<HTMLVideoElement | null>(null)
let player: Plyr | null = null

const initPlayer = () => {
  if (player) {
    player.destroy()
    player = null
  }
  if (videoRef.value) {
    player = new Plyr(videoRef.value, {
      controls: [
        'play-large',
        'play',
        'progress',
        'current-time',
        'mute',
        'volume',
        'settings',
        'fullscreen',
      ],
    })
  }
}

onMounted(() => {
  initPlayer()
})

// Re-initialize when src changes
watch(() => props.src, () => {
  initPlayer()
})

onBeforeUnmount(() => {
  if (player) {
    player.destroy()
    player = null
  }
})
</script>

<template>
  <div id="video-container" class="relative w-full aspect-video bg-black rounded-lg overflow-hidden">
    <video
      ref="videoRef"
      class="w-full h-full"
      :src="src"
      :poster="poster"
      controls
      playsinline
    />
  </div>
</template>

<style>
.plyr--full-ui input[type=range] {
  color: #ef4444;
}
.plyr__control.plyr__control--overlaid {
  background: #ef4444;
  border-radius: 9999px;
}
.plyr--full-ui input[type="range"] {
  color: var(--plyr-range-fill-background, var(--plyr-color-main, #ef4444));
}
</style>