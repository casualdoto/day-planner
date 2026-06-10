export interface DayEntry {
  id: number
  date: string
  planText: string
  resultText: string
  createdAt: string
  updatedAt: string
}

export interface Tag {
  id: number
  name: string
  color: string
}

export interface Task {
  id: number
  title: string
  description: string
  date: string
  dueDate: string | null
  completed: boolean
  completedAt: string | null
  importance: number
  tags: Tag[]
  createdAt: string
  updatedAt: string
}

export interface CreateTaskInput {
  title: string
  description: string
  date: string
  dueDate?: string | null
  importance: number
  tagIds: number[]
}

export interface UpdateTaskInput {
  title?: string
  description?: string
  date?: string
  dueDate?: string | null
  completed?: boolean
  importance?: number
  tagIds?: number[]
}

export interface TaskFilter {
  date?: string | null
  tagId?: number | null
  importance?: number | null
  completed?: boolean | null
  query?: string | null
}

export interface CalendarEvent {
  id: number
  title: string
  description: string
  startDate: string
  endDate: string | null
  completed: boolean
  completedAt: string | null
  createdAt: string
  updatedAt: string
}

export interface CreateCalendarEventInput {
  title: string
  description: string
  startDate: string
  endDate?: string | null
}
