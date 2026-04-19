<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()

const username = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const handleRegister = async () => {
  if (!username.value || !email.value || !password.value) {
    error.value = 'Please fill in all fields'
    return
  }
  error.value = ''
  loading.value = true
  try {
    await axios.post('/api/auth/register', {
      username: username.value,
      email: email.value,
      password: password.value,
    })
    router.push('/login')
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Registration failed. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background">
    <div class="w-full max-w-md space-y-8 p-8">
      <!-- Logo / Brand -->
      <div class="text-center">
        <div class="flex items-center justify-center gap-2 mb-2">
          <div class="w-10 h-10 rounded-full bg-red-500 flex items-center justify-center">
            <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
              <path d="M8 5v14l11-7z"/>
            </svg>
          </div>
          <span class="text-2xl font-bold">StreamVault</span>
        </div>
        <h1 class="text-3xl font-bold tracking-tight">Create account</h1>
        <p class="text-muted-foreground mt-2">Join StreamVault today</p>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleRegister" class="space-y-5">
        <div v-if="error" class="bg-destructive/10 text-destructive text-sm rounded-md px-4 py-3">
          {{ error }}
        </div>

        <div class="space-y-1.5">
          <label for="username" class="text-sm font-medium">Username</label>
          <input
            id="username"
            v-model="username"
            type="text"
            required
            autocomplete="username"
            placeholder="yourname"
            class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
          />
        </div>

        <div class="space-y-1.5">
          <label for="email" class="text-sm font-medium">Email</label>
          <input
            id="email"
            v-model="email"
            type="email"
            required
            autocomplete="email"
            placeholder="you@example.com"
            class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
          />
        </div>

        <div class="space-y-1.5">
          <label for="password" class="text-sm font-medium">Password</label>
          <input
            id="password"
            v-model="password"
            type="password"
            required
            autocomplete="new-password"
            placeholder="••••••••"
            class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
          />
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full rounded-md bg-red-500 px-4 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 disabled:opacity-50 transition-colors"
        >
          <span v-if="loading">Creating account…</span>
          <span v-else>Create account</span>
        </button>
      </form>

      <p class="text-center text-sm text-muted-foreground">
        Already have an account?
        <RouterLink to="/login" class="font-medium text-red-500 hover:underline">Sign in</RouterLink>
      </p>
    </div>
  </div>
</template>
