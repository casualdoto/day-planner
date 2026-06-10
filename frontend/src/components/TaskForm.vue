<script setup lang="ts">
import { ref } from 'vue'
import type { Tag } from '../api/types'

defineProps<{ date: string; tags: Tag[] }>()
const emit = defineEmits<{ submit: [input: { title: string; importance: number; tagIds: number[] }] }>()

const title = ref('')
const importance = ref(2)
const tagId = ref<number | null>(null)

function submit() {
  const cleanTitle = title.value.trim()
  if (!cleanTitle) return
  emit('submit', { title: cleanTitle, importance: importance.value, tagIds: tagId.value ? [tagId.value] : [] })
  title.value = ''
  importance.value = 2
  tagId.value = null
}
</script>

<template>
  <form class="task-form" @submit.prevent="submit">
    <input v-model="title" type="text" placeholder="Quick task" />
    <select v-model.number="importance">
      <option :value="1">Low</option>
      <option :value="2">Medium</option>
      <option :value="3">High</option>
      <option :value="4">Critical</option>
    </select>
    <select v-model="tagId">
      <option :value="null">No tag</option>
      <option v-for="tag in tags" :key="tag.id" :value="tag.id">{{ tag.name }}</option>
    </select>
    <button type="submit">Add</button>
  </form>
</template>
