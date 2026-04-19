<script setup lang="ts">
import VideoPlayer from '@/components/videoPlayer.vue'
import axios from 'axios'
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const videoId = computed(() => route.params.videoId as string)

interface Video {
  id: number
  title: string
  description: string
  video_path: string
  thumbnail_path: string
  created_at: string
}

const video = ref<Video | null>(null)
const loading = ref(true)
const error = ref('')

const streamUrl = computed(() =>
  videoId.value ? `/api/stream/${videoId.value}` : ''
)

const thumbnailUrl = computed(() =>
  videoId.value ? `/api/thumbnail/${videoId.value}` : ''
)

onMounted(async () => {
  if (!videoId.value) return
  try {
    const res = await axios.get(`/api/videos/${videoId.value}`)
    video.value = res.data
  } catch (err: any) {
    error.value = 'Video not found'
    console.error(err)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mt-4">
    <div class="lg:col-span-2 space-y-4">
      <!-- Loading -->
      <div v-if="loading" class="aspect-video bg-muted rounded-lg flex items-center justify-center">
        <div class="animate-spin rounded-full h-10 w-10 border-4 border-red-500 border-t-transparent"></div>
      </div>
      <!-- Error -->
      <div v-else-if="error" class="aspect-video bg-muted rounded-lg flex items-center justify-center">
        <p class="text-muted-foreground">{{ error }}</p>
      </div>
      <!-- Player -->
      <template v-else-if="video">
        <VideoPlayer :src="streamUrl" :poster="thumbnailUrl" />
        <div>
          <h1 class="text-xl font-bold mt-2">{{ video.title }}</h1>
          <div class="mt-3 bg-muted/50 p-3 rounded-lg">
            <p class="text-sm text-muted-foreground">
              Uploaded {{ new Date(video.created_at).toLocaleDateString() }}
            </p>
            <p class="mt-2 whitespace-pre-line text-sm">{{ video.description }}</p>
          </div>
        </div>
      </template>
    </div>

    <!-- Related sidebar placeholder -->
    <div class="space-y-4">
      <h3 class="font-medium">Related videos</h3>
      <p class="text-sm text-muted-foreground">More coming soon...</p>
    </div>
  </div>
</template>