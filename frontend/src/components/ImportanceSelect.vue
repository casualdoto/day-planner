<script setup lang="ts">
defineProps<{ modelValue: number | null }>()
defineEmits<{ 'update:modelValue': [value: number | null] }>()
</script>

<template>
  <select :value="modelValue ?? ''" @change="$emit('update:modelValue', ($event.target as HTMLSelectElement).value ? Number(($event.target as HTMLSelectElement).value) : null)">
    <option value="">Any importance</option>
    <option :value="1">Low</option>
    <option :value="2">Medium</option>
    <option :value="3">High</option>
    <option :value="4">Critical</option>
  </select>
</template>
