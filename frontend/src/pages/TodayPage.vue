<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import DailyTextEditor from '../components/DailyTextEditor.vue'
import TaskForm from '../components/TaskForm.vue'
import TaskList from '../components/TaskList.vue'
import { getDayEntry, saveDayPlan, saveDayResult } from '../api/day'
import { createTask, deleteTask, listTasks, toggleTaskCompleted } from '../api/tasks'
import { listTags } from '../api/tags'
import type { DayEntry, Tag, Task } from '../api/types'
import { toDateOnly } from '../utils/date'

const selectedDate = ref(toDateOnly(new Date()))
const day = ref<DayEntry | null>(null)
const tasks = ref<Task[]>([])
const tags = ref<Tag[]>([])
const planStatus = ref('')
const resultStatus = ref('')
const error = ref('')
let planTimer = 0
let resultTimer = 0

async function load() {
  try {
    error.value = ''
    day.value = await getDayEntry(selectedDate.value)
    tags.value = await listTags()
    tasks.value = await listTasks({ date: selectedDate.value })
  } catch (err) {
    error.value = err instanceof Error ? err.message : String(err)
  }
}

async function addTask(input: { title: string; importance: number; tagIds: number[] }) {
  await createTask({ ...input, date: selectedDate.value, description: '' })
  await load()
}

async function removeTask(id: number) {
  if (!confirm('Delete this task?')) return
  await deleteTask(id)
  await load()
}

watch(selectedDate, load)
watch(() => day.value?.planText, (value, oldValue) => {
  if (oldValue === undefined || !day.value) return
  clearTimeout(planTimer)
  planStatus.value = 'Saving...'
  planTimer = window.setTimeout(async () => {
    try {
      await saveDayPlan(selectedDate.value, value ?? '')
      planStatus.value = 'Saved'
    } catch (err) {
      planStatus.value = 'Save failed'
      error.value = err instanceof Error ? err.message : String(err)
    }
  }, 700)
})
watch(() => day.value?.resultText, (value, oldValue) => {
  if (oldValue === undefined || !day.value) return
  clearTimeout(resultTimer)
  resultStatus.value = 'Saving...'
  resultTimer = window.setTimeout(async () => {
    try {
      await saveDayResult(selectedDate.value, value ?? '')
      resultStatus.value = 'Saved'
    } catch (err) {
      resultStatus.value = 'Save failed'
      error.value = err instanceof Error ? err.message : String(err)
    }
  }, 700)
})

onMounted(load)
</script>

<template>
  <section class="page">
    <header class="page-header">
      <div>
        <p class="eyebrow">Notebook</p>
        <h2>Today</h2>
      </div>
      <input v-model="selectedDate" type="date" />
    </header>

    <p v-if="error" class="error">{{ error }}</p>

    <div v-if="day" class="today-layout">
      <DailyTextEditor v-model="day.planText" title="Plan" placeholder="Write the shape of the day..." :save-status="planStatus" />
      <section class="tasks-panel">
        <h2>Tasks</h2>
        <TaskForm :date="selectedDate" :tags="tags" @submit="addTask" />
        <TaskList :tasks="tasks" @toggle="async (id) => { await toggleTaskCompleted(id); await load() }" @delete="removeTask" />
      </section>
      <DailyTextEditor v-model="day.resultText" title="Results" placeholder="What happened, what changed, what matters?" :save-status="resultStatus" />
    </div>
  </section>
</template>
