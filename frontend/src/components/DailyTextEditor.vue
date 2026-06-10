<script setup lang="ts">
defineProps<{ title: string; modelValue: string; placeholder: string; saveStatus?: string }>()
defineEmits<{ 'update:modelValue': [value: string] }>()
</script>

<template>
  <section class="editor-block">
    <div class="editor-header">
      <h2>{{ title }}</h2>
      <span v-if="saveStatus" class="save-status">{{ saveStatus }}</span>
    </div>
    <textarea
      :value="modelValue"
      :placeholder="placeholder"
      @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
    />
  </section>
</template>
