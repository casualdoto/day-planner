<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import CalendarMonth from '../components/CalendarMonth.vue'
import DaySidePanel from '../components/DaySidePanel.vue'
import { getDayEntry } from '../api/day'
import { listCalendarEvents } from '../api/calendar'
import { listTasks } from '../api/tasks'
import type { CalendarEvent, DayEntry, Task } from '../api/types'
import { toDateOnly } from '../utils/date'

const selectedDate = ref(toDateOnly(new Date()))
const month = ref(new Date())
const tasks = ref<Task[]>([])
const events = ref<CalendarEvent[]>([])
const day = ref<DayEntry | null>(null)
const error = ref('')
const monthLabel = computed(() => month.value.toLocaleDateString('en-US', { month: 'long', year: 'numeric' }))

function monthRange() {
  const start = new Date(month.value.getFullYear(), month.value.getMonth(), 1)
  const end = new Date(month.value.getFullYear(), month.value.getMonth() + 1, 0)
  return [toDateOnly(start), toDateOnly(end)]
}

async function load() {
  try {
    error.value = ''
    const [start, end] = monthRange()
    tasks.value = await listTasks({})
    events.value = await listCalendarEvents(start, end)
    day.value = await getDayEntry(selectedDate.value)
  } catch (err) {
    error.value = err instanceof Error ? err.message : String(err)
  }
}

function shiftMonth(delta: number) {
  month.value = new Date(month.value.getFullYear(), month.value.getMonth() + delta, 1)
}

watch([month, selectedDate], load)
onMounted(load)
</script>

<template>
  <section class="page calendar-page">
    <header class="page-header">
      <div>
        <p class="eyebrow">Calendar</p>
        <h2>{{ monthLabel }}</h2>
      </div>
      <div class="button-row">
        <button type="button" @click="shiftMonth(-1)">Previous</button>
        <button type="button" @click="shiftMonth(1)">Next</button>
      </div>
    </header>
    <p v-if="error" class="error">{{ error }}</p>
    <div class="calendar-layout">
      <CalendarMonth :month="month" :selected-date="selectedDate" :tasks="tasks" :events="events" @select="selectedDate = $event" />
      <DaySidePanel :date="selectedDate" :day="day" :tasks="tasks.filter((task) => task.date === selectedDate)" :events="events.filter((event) => event.startDate <= selectedDate && (event.endDate ?? event.startDate) >= selectedDate)" />
    </div>
  </section>
</template>
