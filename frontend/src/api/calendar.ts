import { bridge } from './bridge'
import type { CreateCalendarEventInput } from './types'

export const createCalendarEvent = (input: CreateCalendarEventInput) => bridge().createCalendarEvent(input)
export const updateCalendarEvent = (id: number, input: CreateCalendarEventInput) => bridge().updateCalendarEvent(id, input)
export const deleteCalendarEvent = (id: number) => bridge().deleteCalendarEvent(id)
export const toggleCalendarEventCompleted = (id: number) => bridge().toggleCalendarEventCompleted(id)
export const listCalendarEvents = (startDate: string, endDate: string) => bridge().listCalendarEvents(startDate, endDate)
