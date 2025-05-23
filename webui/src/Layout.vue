<script setup lang="ts">
import Navbar from './components/Navbar.vue';
import Sidebar from './components/Sidebar.vue';

</script>

<template>
      <Navbar />
       <div className="flex">
            <Sidebar />
            <main className="flex-1 md:ml-60">
              <div className="container mx-auto px-4 py-4 max-w-7xl">
                <slot />
              </div>
            </main>
          </div>
</template>