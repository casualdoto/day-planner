import { bridge } from './bridge'

export const getDayEntry = (date: string) => bridge().getDayEntry(date)
export const saveDayPlan = (date: string, planText: string) => bridge().saveDayPlan(date, planText)
export const saveDayResult = (date: string, resultText: string) => bridge().saveDayResult(date, resultText)
