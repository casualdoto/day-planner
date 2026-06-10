<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { createTag, deleteTag, listTags, updateTag } from '../api/tags'
import type { Tag } from '../api/types'

const tags = ref<Tag[]>([])
const name = ref('')
const color = ref('#64748b')
const error = ref('')

async function load() {
  tags.value = await listTags()
}

async function addTag() {
  try {
    error.value = ''
    await createTag(name.value, color.value)
    name.value = ''
    color.value = '#64748b'
    await load()
  } catch (err) {
    error.value = err instanceof Error ? err.message : String(err)
  }
}

async function removeTag(id: number) {
  if (!confirm('Delete this tag?')) return
  await deleteTag(id)
  await load()
}

async function saveTag(tag: Tag) {
  try {
    error.value = ''
    await updateTag(tag.id, tag.name, tag.color)
    await load()
  } catch (err) {
    error.value = err instanceof Error ? err.message : String(err)
  }
}

onMounted(load)
</script>

<template>
  <section class="page">
    <header class="page-header">
      <div>
        <p class="eyebrow">Labels</p>
        <h2>Tags</h2>
      </div>
    </header>
    <p v-if="error" class="error">{{ error }}</p>
    <form class="task-form tag-form" @submit.prevent="addTag">
      <input v-model="name" placeholder="Tag name" />
      <label class="color-control" title="Choose tag color">
        <span class="color-swatch" :style="{ backgroundColor: color }"></span>
        <input v-model="color" type="color" />
      </label>
      <button type="submit">Create</button>
    </form>
    <ul class="tag-list">
      <li v-for="tag in tags" :key="tag.id">
        <input v-model="tag.name" />
        <label class="color-control" title="Choose tag color">
          <span class="color-swatch" :style="{ backgroundColor: tag.color }"></span>
          <input v-model="tag.color" type="color" />
        </label>
        <input v-model="tag.color" class="color-text" />
        <button type="button" class="ghost" @click="saveTag(tag)">Save</button>
        <button type="button" class="ghost danger" @click="removeTag(tag.id)">Delete</button>
      </li>
    </ul>
  </section>
</template>
