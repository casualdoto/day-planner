import { bridge } from './bridge'
import type { CreateTaskInput, TaskFilter, UpdateTaskInput } from './types'

export const createTask = (input: CreateTaskInput) => bridge().createTask(input)
export const updateTask = (id: number, input: UpdateTaskInput) => bridge().updateTask(id, input)
export const deleteTask = (id: number) => bridge().deleteTask(id)
export const toggleTaskCompleted = (id: number) => bridge().toggleTaskCompleted(id)
export const listTasks = (filter: TaskFilter) => bridge().listTasks(filter)
