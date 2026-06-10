<script setup lang="ts">
import type { Task } from '../api/types'

defineProps<{ task: Task }>()
defineEmits<{ toggle: [id: number]; delete: [id: number] }>()
</script>

<template>
  <li class="task-item" :class="{ completed: task.completed }">
    <label>
      <input type="checkbox" :checked="task.completed" @change="$emit('toggle', task.id)" />
      <span>{{ task.title }}</span>
    </label>
    <div class="task-meta">
      <span>Importance {{ task.importance }}</span>
      <span v-for="tag in task.tags" :key="tag.id" class="tag-pill" :style="{ backgroundColor: tag.color }">{{ tag.name }}</span>
      <button type="button" class="ghost danger" @click="$emit('delete', task.id)">Delete</button>
    </div>
  </li>
</template>
