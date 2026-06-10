<script setup lang="ts">
import type { Task } from '../api/types'
import TaskItem from './TaskItem.vue'

defineProps<{ tasks: Task[] }>()
defineEmits<{ toggle: [id: number]; delete: [id: number] }>()
</script>

<template>
  <ul v-if="tasks.length" class="task-list">
    <TaskItem v-for="task in tasks" :key="task.id" :task="task" @toggle="$emit('toggle', $event)" @delete="$emit('delete', $event)" />
  </ul>
  <p v-else class="empty">No tasks here yet.</p>
</template>
