<script setup lang="ts">
import { computed, ref } from 'vue'
import AppShell from './components/AppShell.vue'
import TodayPage from './pages/TodayPage.vue'
import TasksPage from './pages/TasksPage.vue'
import CalendarPage from './pages/CalendarPage.vue'
import TagsPage from './pages/TagsPage.vue'

const page = ref('today')
const pages = { today: TodayPage, tasks: TasksPage, calendar: CalendarPage, tags: TagsPage }
const currentPage = computed(() => pages[page.value as keyof typeof pages] ?? TodayPage)
</script>

<template>
  <AppShell :active-page="page" @navigate="page = $event">
    <component :is="currentPage" />
  </AppShell>
</template>
