<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import { Upload, X, CheckCircle, Loader } from 'lucide-vue-next'

const router = useRouter()

const CHUNK_SIZE = 5 << 20 // 5 MB

// State
const file = ref<File | null>(null)
const title = ref('')
const description = ref('')
const isDragging = ref(false)
const uploading = ref(false)
const progress = ref(0)
const uploadDone = ref(false)
const error = ref('')

// File handling
const handleFileChange = (e: Event) => {
  const input = e.target as HTMLInputElement
  if (input.files?.length) selectFile(input.files[0])
}

const handleDrop = (e: DragEvent) => {
  e.preventDefault()
  isDragging.value = false
  const dropped = e.dataTransfer?.files?.[0]
  if (dropped?.type.startsWith('video/')) selectFile(dropped)
  else error.value = 'Please select a valid video file'
}

const selectFile = (f: File) => {
  file.value = f
  if (!title.value) title.value = f.name.replace(/\.[^.]+$/, '')
  error.value = ''
}

const clearFile = () => {
  file.value = null
  progress.value = 0
  uploadDone.value = false
  error.value = ''
}

// Format file size
const formatSize = (bytes: number) => {
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(2)} GB`
}

// Chunked upload
const submit = async () => {
  if (!file.value) { error.value = 'Please select a video file'; return }
  if (!title.value.trim()) { error.value = 'Please enter a title'; return }

  uploading.value = true
  error.value = ''
  progress.value = 0

  const f = file.value
  const totalChunks = Math.ceil(f.size / CHUNK_SIZE)

  try {
    for (let i = 0; i < totalChunks; i++) {
      const start = i * CHUNK_SIZE
      const end = Math.min(start + CHUNK_SIZE, f.size)
      const chunk = f.slice(start, end)

      const form = new FormData()
      form.append('file', chunk)
      form.append('file_name', f.name)
      form.append('index', String(i))
      form.append('total', String(totalChunks))

      await axios.post('/api/videos/upload', form)
      progress.value = Math.round(((i + 1) / totalChunks) * 100)
    }

    uploadDone.value = true
    // wait a moment then redirect home
    setTimeout(() => router.push('/'), 1500)
  } catch (err: any) {
    console.log(err);

    error.value = err.response?.data?.error || 'Upload failed. Please try again.'
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <div class="max-w-2xl mx-auto py-8 px-4">
    <h1 class="text-2xl font-bold mb-6">Upload Video</h1>

    <!-- Success state -->
    <div v-if="uploadDone" class="flex flex-col items-center justify-center py-16 gap-4 text-center">
      <CheckCircle class="w-16 h-16 text-green-500" />
      <h2 class="text-xl font-semibold">Upload successful!</h2>
      <p class="text-muted-foreground">Redirecting to home…</p>
    </div>

    <template v-else>
      <!-- Drop Zone (when no file selected) -->
      <div v-if="!file" class="border-2 border-dashed rounded-xl p-16 text-center transition-colors cursor-pointer"
        :class="isDragging ? 'border-red-500 bg-red-50 dark:bg-red-950/20' : 'border-muted-foreground/30 hover:border-red-400'"
        @dragover.prevent="isDragging = true" @dragleave="isDragging = false" @drop="handleDrop"
        @click="($refs.fileInput as HTMLInputElement).click()">
        <Upload class="w-12 h-12 mx-auto mb-4 text-muted-foreground" />
        <h2 class="text-lg font-medium mb-1">Drag and drop a video</h2>
        <p class="text-sm text-muted-foreground mb-4">MP4, WebM, MOV — up to 2 GB</p>
        <button type="button"
          class="px-4 py-2 rounded-md bg-red-500 text-white text-sm font-medium hover:bg-red-600 transition-colors">
          Select file
        </button>
        <input ref="fileInput" type="file" accept="video/*" class="hidden" @change="handleFileChange" />
      </div>

      <!-- Form (file selected) -->
      <form v-else @submit.prevent="submit" class="space-y-5">
        <!-- File info row -->
        <div class="flex items-center gap-3 p-4 rounded-lg border bg-muted/30">
          <div class="flex-1 min-w-0">
            <p class="font-medium truncate">{{ file.name }}</p>
            <p class="text-sm text-muted-foreground">{{ formatSize(file.size) }}</p>
          </div>
          <button type="button" @click="clearFile" :disabled="uploading"
            class="p-1.5 rounded-md hover:bg-muted transition-colors text-muted-foreground disabled:opacity-50">
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Title -->
        <div class="space-y-1.5">
          <label for="title" class="text-sm font-medium">
            Title <span class="text-red-500">*</span>
          </label>
          <input id="title" v-model="title" type="text" required maxlength="100"
            placeholder="Enter a title for your video"
            class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-red-500" />
          <p class="text-xs text-muted-foreground text-right">{{ title.length }}/100</p>
        </div>

        <!-- Description -->
        <div class="space-y-1.5">
          <label for="description" class="text-sm font-medium">Description</label>
          <textarea id="description" v-model="description" rows="4" maxlength="2000"
            placeholder="Tell viewers about your video"
            class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-red-500 resize-y" />
          <p class="text-xs text-muted-foreground text-right">{{ description.length }}/2000</p>
        </div>

        <!-- Error -->
        <div v-if="error" class="rounded-md bg-destructive/10 text-destructive px-4 py-3 text-sm">
          {{ error }}
        </div>

        <!-- Progress -->
        <div v-if="uploading" class="space-y-1.5">
          <div class="flex justify-between text-sm text-muted-foreground">
            <span>Uploading…</span>
            <span>{{ progress }}%</span>
          </div>
          <div class="w-full h-2 bg-muted rounded-full overflow-hidden">
            <div class="h-full bg-red-500 rounded-full transition-all duration-300"
              :style="{ width: progress + '%' }" />
          </div>
        </div>

        <!-- Submit -->
        <button type="submit" :disabled="uploading"
          class="w-full flex items-center justify-center gap-2 rounded-md bg-red-500 px-4 py-2.5 text-sm font-semibold text-white hover:bg-red-600 disabled:opacity-60 transition-colors">
          <Loader v-if="uploading" class="w-4 h-4 animate-spin" />
          <Upload v-else class="w-4 h-4" />
          <span>{{ uploading ? 'Uploading…' : 'Upload Video' }}</span>
        </button>
      </form>
    </template>
  </div>
</template>