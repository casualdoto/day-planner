import type {
  CalendarEvent,
  CreateCalendarEventInput,
  CreateTaskInput,
  DayEntry,
  Tag,
  Task,
  TaskFilter,
  UpdateTaskInput,
} from './types'

type Bridge = {
  getDayEntry(date: string): Promise<DayEntry>
  saveDayPlan(date: string, planText: string): Promise<void>
  saveDayResult(date: string, resultText: string): Promise<void>
  createTask(input: CreateTaskInput): Promise<Task>
  updateTask(id: number, input: UpdateTaskInput): Promise<Task>
  deleteTask(id: number): Promise<void>
  toggleTaskCompleted(id: number): Promise<Task>
  listTasks(filter: TaskFilter): Promise<Task[]>
  createTag(name: string, color: string): Promise<Tag>
  updateTag(id: number, name: string, color: string): Promise<Tag>
  deleteTag(id: number): Promise<void>
  listTags(): Promise<Tag[]>
  createCalendarEvent(input: CreateCalendarEventInput): Promise<CalendarEvent>
  updateCalendarEvent(id: number, input: CreateCalendarEventInput): Promise<CalendarEvent>
  deleteCalendarEvent(id: number): Promise<void>
  toggleCalendarEventCompleted(id: number): Promise<CalendarEvent>
  listCalendarEvents(startDate: string, endDate: string): Promise<CalendarEvent[]>
}

const devApiBase = `${window.location.protocol}//${window.location.hostname}:8080`
const apiBase = import.meta.env.VITE_API_BASE_URL || (window.location.port === '5173' ? devApiBase : '')

export function bridge(): Bridge {
  return httpBridge
}

const httpBridge: Bridge = {
  getDayEntry(date) {
    return request(`/api/day?date=${encodeURIComponent(date)}`)
  },
  saveDayPlan(date, planText) {
    return request('/api/day/plan', { date, planText })
  },
  saveDayResult(date, resultText) {
    return request('/api/day/result', { date, resultText })
  },
  createTask(input) {
    return request('/api/tasks/create', input)
  },
  updateTask(id, input) {
    return request('/api/tasks/update', { id, input })
  },
  deleteTask(id) {
    return request('/api/tasks/delete', { id })
  },
  toggleTaskCompleted(id) {
    return request('/api/tasks/toggle', { id })
  },
  listTasks(filter) {
    return request<Task[] | null>('/api/tasks/list', filter).then((tasks) => tasks ?? [])
  },
  createTag(name, color) {
    return request('/api/tags/create', { name, color })
  },
  updateTag(id, name, color) {
    return request('/api/tags/update', { id, name, color })
  },
  deleteTag(id) {
    return request('/api/tags/delete', { id })
  },
  listTags() {
    return request<Tag[] | null>('/api/tags').then((tags) => tags ?? [])
  },
  createCalendarEvent(input) {
    return request('/api/calendar/create', input)
  },
  updateCalendarEvent(id, input) {
    return request('/api/calendar/update', { id, input })
  },
  deleteCalendarEvent(id) {
    return request('/api/calendar/delete', { id })
  },
  toggleCalendarEventCompleted(id) {
    return request('/api/calendar/toggle', { id })
  },
  listCalendarEvents(startDate, endDate) {
    return request<CalendarEvent[] | null>(`/api/calendar?startDate=${encodeURIComponent(startDate)}&endDate=${encodeURIComponent(endDate)}`).then((events) => events ?? [])
  },
}

async function request<T>(path: string, body?: unknown): Promise<T> {
  const response = await fetch(`${apiBase}${path}`, {
    method: body === undefined ? 'GET' : 'POST',
    headers: body === undefined ? undefined : { 'Content-Type': 'application/json' },
    body: body === undefined ? undefined : JSON.stringify(body),
  })
  if (!response.ok) {
    const message = await response.json().catch(() => null)
    throw new Error(message?.error ?? `Request failed with ${response.status}`)
  }
  if (response.status === 204) {
    return undefined as T
  }
  return response.json()
}
