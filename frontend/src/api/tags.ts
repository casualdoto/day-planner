import { bridge } from './bridge'

export const createTag = (name: string, color: string) => bridge().createTag(name, color)
export const updateTag = (id: number, name: string, color: string) => bridge().updateTag(id, name, color)
export const deleteTag = (id: number) => bridge().deleteTag(id)
export const listTags = () => bridge().listTags()
