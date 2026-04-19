<script setup lang="ts">
defineProps<{
  video: {
    id: number
    title: string
    description?: string
    video_path?: string
    thumbnail_path?: string
    created_at?: string
    user_id?: number
  } | null
  horizontal?: boolean
}>()
</script>

<template>
  <!-- Horizontal layout (sidebar related videos) -->
  <RouterLink
    v-if="horizontal && video"
    :to="`/watch/${video.id}`"
    class="group flex gap-3 hover:bg-muted/50 rounded-lg p-2"
  >
    <div class="relative flex-shrink-0 w-40 h-24 rounded-lg overflow-hidden bg-muted">
      <img
        v-if="video.thumbnail_path"
        :src="`/api/thumbnail/${video.id}`"
        :alt="video.title"
        class="w-full h-full object-cover"
      />
      <div v-else class="w-full h-full flex items-center justify-center text-muted-foreground">
        <svg class="w-8 h-8" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
      </div>
    </div>
    <div class="flex flex-col flex-1 min-w-0">
      <h3 class="font-medium text-sm line-clamp-2 group-hover:text-red-500">{{ video.title }}</h3>
      <p class="text-xs text-muted-foreground mt-1" v-if="video.created_at">
        {{ new Date(video.created_at).toLocaleDateString() }}
      </p>
    </div>
  </RouterLink>

  <!-- Vertical / grid layout (home page) -->
  <RouterLink
    v-else-if="video"
    :to="`/watch/${video.id}`"
    class="group block"
  >
    <div class="relative aspect-video rounded-lg overflow-hidden mb-2 bg-muted">
      <img
        v-if="video.thumbnail_path"
        :src="`/api/thumbnail/${video.id}`"
        :alt="video.title"
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-200"
      />
      <div v-else class="w-full h-full flex items-center justify-center text-muted-foreground group-hover:scale-105 transition-transform duration-200">
        <svg class="w-12 h-12" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
      </div>
    </div>
    <div>
      <h3 class="font-medium text-sm line-clamp-2 group-hover:text-red-500">{{ video.title }}</h3>
      <p class="text-xs text-muted-foreground mt-1" v-if="video.created_at">
        {{ new Date(video.created_at).toLocaleDateString() }}
      </p>
    </div>
  </RouterLink>

  <!-- Skeleton placeholder (no video) -->
  <div v-else class="animate-pulse">
    <div class="aspect-video rounded-lg bg-muted mb-2"></div>
    <div class="h-4 bg-muted rounded w-3/4 mb-1"></div>
    <div class="h-3 bg-muted rounded w-1/2"></div>
  </div>
</template>