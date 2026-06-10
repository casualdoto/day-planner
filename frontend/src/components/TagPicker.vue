<script setup lang="ts">
import type { Tag } from '../api/types'

const props = defineProps<{ tags: Tag[]; modelValue: number | null }>()
const emit = defineEmits<{ 'update:modelValue': [value: number | null] }>()

function update(value: string) {
  emit('update:modelValue', value ? Number(value) : null)
}
</script>

<template>
  <select :value="props.modelValue ?? ''" @change="update(($event.target as HTMLSelectElement).value)">
    <option value="">Any tag</option>
    <option v-for="tag in tags" :key="tag.id" :value="tag.id">{{ tag.name }}</option>
  </select>
</template>
