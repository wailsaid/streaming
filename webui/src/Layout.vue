<script setup lang="ts">
import { useRoute } from 'vue-router';
import Navbar from './components/Navbar.vue';
import Sidebar from './components/Sidebar.vue';

const useroute = useRoute();

</script>

<template>
  <!--    <Navbar />
       <div className="flex">
            <Sidebar v-if="useroute.name !== 'watch'" />
            <main class="flex-1 md:ml-60">
              <div class="container mx-auto px-4 py-4 max-w-7xl">
                <slot /> 
              </div>
            </main>
          </div> -->

  <Navbar />
  <template v-if="useroute.name === 'watch'">
    <main class="container mx-auto px-4 py-4 max-w-7xl">
      <slot />
    </main>

  </template>
  <template v-else>
    <div className="flex">
      <Sidebar v-if="useroute.name !== 'watch'" />
      <main class="flex-1 md:ml-60">
        <div class="container mx-auto px-4 py-4 max-w-7xl">
          <slot />
        </div>
      </main>
    </div>
  </template>
</template>