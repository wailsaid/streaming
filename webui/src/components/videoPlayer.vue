<script setup >
import Plyr from 'plyr';
import 'plyr/dist/plyr.css'
import { onBeforeUnmount, onMounted, ref } from 'vue';

const videoRef = ref(null)
let player = null

onMounted(()=>{
    player = new Plyr(videoRef.value,{
    controls: [
      'play-large',
      'play',
      'progress',
      'current-time',
      'mute',
      'volume',
      'settings',
      'fullscreen'
    ],
  })
})

onBeforeUnmount(()=>{
    if (player) {
      player.destroy();
      player = null;
        
    }
})

</script>
<template>
    <div id="video-container" class="relative w-full aspect-video bg-black rounded-lg overflow-hidden">
      <video
        ref="videoRef"
        class="w-full h-full"
        src="https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4"
        poster="data:image/jpeg;base64,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"
        controls
        playsinline
      />
    </div>
</template>


<style>
/* Optional: Customize Plyr theme colors (overwrite default) */
.plyr--full-ui input[type=range] {
  color: #aa4e4e; /* Tailwind's blue-500 */
}

.plyr__control.plyr__control--overlaid {
  background: #aa4e4e;
  border-radius: 9999px;
}

.plyr--full-ui input[type="range"]{

    color: var(--plyr-range-fill-background,var(--plyr-color-main,var(--plyr-color-main,#aa4e4e)));
}
</style>