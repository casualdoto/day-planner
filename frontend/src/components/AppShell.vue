<script setup lang="ts">
import Sidebar from './Sidebar.vue'

defineProps<{ activePage: string }>()
defineEmits<{ navigate: [page: string] }>()
</script>

<template>
  <div class="app-shell">
    <Sidebar :active-page="activePage" @navigate="$emit('navigate', $event)" />
    <main class="main-panel">
      <slot />
    </main>
  </div>
</template>
