<script setup lang="ts">
import { computed } from 'vue'
import type { CalendarEvent, Task } from '../api/types'
import { toDateOnly } from '../utils/date'

const props = defineProps<{ month: Date; selectedDate: string; tasks: Task[]; events: CalendarEvent[] }>()
defineEmits<{ select: [date: string] }>()

const today = toDateOnly(new Date())
const days = computed(() => {
  const year = props.month.getFullYear()
  const month = props.month.getMonth()
  const first = new Date(year, month, 1)
  const start = new Date(first)
  start.setDate(first.getDate() - first.getDay())
  return Array.from({ length: 42 }, (_, index) => {
    const date = new Date(start)
    date.setDate(start.getDate() + index)
    const iso = toDateOnly(date)
    const dayTasks = props.tasks.filter((task) => task.date === iso)
    return {
      iso,
      label: date.getDate(),
      currentMonth: date.getMonth() === month,
      today: iso === today,
      selected: iso === props.selectedDate,
      hasTasks: dayTasks.length > 0,
      hasCompleted: dayTasks.some((task) => task.completed),
      hasImportant: dayTasks.some((task) => task.importance >= 3),
      hasEvents: props.events.some((event) => event.startDate <= iso && (event.endDate ?? event.startDate) >= iso),
    }
  })
})
</script>

<template>
  <div class="calendar-grid">
    <span v-for="dayName in ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']" :key="dayName" class="weekday">{{ dayName }}</span>
    <button
      v-for="day in days"
      :key="day.iso"
      type="button"
      class="calendar-day"
      :class="{ muted: !day.currentMonth, today: day.today, selected: day.selected }"
      @click="$emit('select', day.iso)"
    >
      <span>{{ day.label }}</span>
      <i v-if="day.hasTasks" class="dot task-dot"></i>
      <i v-if="day.hasCompleted" class="dot done-dot"></i>
      <i v-if="day.hasImportant" class="dot important-dot"></i>
      <i v-if="day.hasEvents" class="dot event-dot"></i>
    </button>
  </div>
</template>
