<script setup lang="ts">
import VideoCard from '@/components/videoCard.vue'
import axios from 'axios'
import { onMounted, ref } from 'vue'

interface Video {
  id: number
  title: string
  description: string
  video_path: string
  thumbnail_path: string
  user_id: number
  created_at: string
}

const videos = ref<Video[]>([])
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    const res = await axios.get('/api/videos')
    videos.value = res.data ?? []
  } catch (err: any) {
    error.value = 'Could not load videos'
    console.error('Error fetching videos:', err)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="home">
    <div v-if="loading" class="flex items-center justify-center py-20">
      <div class="animate-spin rounded-full h-10 w-10 border-4 border-red-500 border-t-transparent"></div>
    </div>

    <div v-else-if="error" class="text-center py-20 text-muted-foreground">
      {{ error }}
    </div>

    <div v-else-if="videos.length === 0" class="text-center py-20 text-muted-foreground">
      <p class="text-lg">No videos yet.</p>
      <p class="text-sm mt-1">
        <RouterLink to="/profile/upload" class="text-red-500 hover:underline">Upload your first video</RouterLink>
      </p>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mt-4">
      <VideoCard
        v-for="video in videos"
        :key="video.id"
        :video="video"
      />
    </div>
  </div>
</template>