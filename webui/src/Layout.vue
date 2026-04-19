<script setup lang="ts">
import { useRoute } from 'vue-router';
import Navbar from './components/Navbar.vue';
import Sidebar from './components/Sidebar.vue';

const route = useRoute();

// Don't show app shell on auth pages
const isAuthPage = () => ['login', 'register'].includes(route.name as string)
</script>

<template>
  <!-- Auth pages: no navbar or sidebar -->
  <template v-if="isAuthPage()">
    <slot />
  </template>

  <!-- App pages: full layout -->
  <template v-else>
    <Navbar />
    <div class="flex">
      <Sidebar v-if="route.name !== 'watch'" />
      <main :class="['flex-1', route.name !== 'watch' ? 'md:ml-60' : '']">
        <div class="container mx-auto px-4 py-4 max-w-7xl">
          <slot />
        </div>
      </main>
    </div>
  </template>
</template>