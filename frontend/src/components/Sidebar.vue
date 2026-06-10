<script setup lang="ts">
defineProps<{ activePage: string }>()
defineEmits<{ navigate: [page: string] }>()

const items = [
  ['today', 'Today'],
  ['tasks', 'Tasks'],
  ['calendar', 'Calendar'],
  ['tags', 'Tags'],
]
</script>

<template>
  <aside class="sidebar">
    <h1>Day Planner</h1>
    <nav>
      <button
        v-for="[key, label] in items"
        :key="key"
        :class="{ active: activePage === key }"
        type="button"
        @click="$emit('navigate', key)"
      >
        {{ label }}
      </button>
    </nav>
  </aside>
</template>
