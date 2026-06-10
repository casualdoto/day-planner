<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import ImportanceSelect from '../components/ImportanceSelect.vue'
import TagPicker from '../components/TagPicker.vue'
import TaskForm from '../components/TaskForm.vue'
import TaskList from '../components/TaskList.vue'
import { createTask, deleteTask, listTasks, toggleTaskCompleted } from '../api/tasks'
import { listTags } from '../api/tags'
import type { Tag, Task } from '../api/types'

const today = new Date().toISOString().slice(0, 10)
const tasks = ref<Task[]>([])
const tags = ref<Tag[]>([])
const date = ref('')
const status = ref('')
const tagId = ref<number | null>(null)
const importance = ref<number | null>(null)
const query = ref('')
const error = ref('')

async function load() {
  try {
    error.value = ''
    tags.value = await listTags()
    tasks.value = await listTasks({
      date: date.value || null,
      tagId: tagId.value,
      importance: importance.value,
      completed: status.value === '' ? null : status.value === 'done',
      query: query.value || null,
    })
  } catch (err) {
    error.value = err instanceof Error ? err.message : String(err)
  }
}

async function addTask(input: { title: string; importance: number; tagIds: number[] }) {
  await createTask({ ...input, date: date.value || today, description: '' })
  await load()
}

async function removeTask(id: number) {
  if (!confirm('Delete this task?')) return
  await deleteTask(id)
  await load()
}

watch([date, status, tagId, importance, query], load)
onMounted(load)
</script>

<template>
  <section class="page">
    <header class="page-header">
      <div>
        <p class="eyebrow">Todo</p>
        <h2>Tasks</h2>
      </div>
    </header>
    <p v-if="error" class="error">{{ error }}</p>
    <TaskForm :date="date || today" :tags="tags" @submit="addTask" />
    <div class="filters">
      <input v-model="date" type="date" />
      <select v-model="status">
        <option value="">Any status</option>
        <option value="open">Open</option>
        <option value="done">Done</option>
      </select>
      <TagPicker v-model="tagId" :tags="tags" />
      <ImportanceSelect v-model="importance" />
      <input v-model="query" type="search" placeholder="Search" />
    </div>
    <TaskList :tasks="tasks" @toggle="async (id) => { await toggleTaskCompleted(id); await load() }" @delete="removeTask" />
  </section>
</template>
