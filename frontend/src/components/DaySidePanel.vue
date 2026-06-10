<script setup lang="ts">
import type { CalendarEvent, DayEntry, Task } from '../api/types'

defineProps<{ date: string; day: DayEntry | null; tasks: Task[]; events: CalendarEvent[] }>()
</script>

<template>
  <aside class="side-panel">
    <h2>{{ date }}</h2>
    <h3>Plan</h3>
    <p>{{ day?.planText || 'No plan yet.' }}</p>
    <h3>Tasks</h3>
    <p>{{ tasks.length }} task(s)</p>
    <h3>Events</h3>
    <p>{{ events.length }} event(s)</p>
    <h3>Result</h3>
    <p>{{ day?.resultText || 'No review yet.' }}</p>
  </aside>
</template>
