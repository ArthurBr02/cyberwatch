import axios from 'axios'
import type { Source, Article, CreateSourcePayload } from '../types'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  headers: { 'Content-Type': 'application/json' }
})

// Sources
export const getSources = () => api.get<Source[]>('/api/v1/sources').then(r => r.data)
export const getSource = (id: number) => api.get<Source>(`/api/v1/sources/${id}`).then(r => r.data)
export const createSource = (data: CreateSourcePayload) => api.post<Source>('/api/v1/sources', data).then(r => r.data)
export const updateSource = (id: number, data: Partial<Source>) => api.put<Source>(`/api/v1/sources/${id}`, data).then(r => r.data)
export const deleteSource = (id: number) => api.delete(`/api/v1/sources/${id}`)
export const generateRules = (id: number) => api.post<Source>(`/api/v1/sources/${id}/generate-rules`).then(r => r.data)
export const testScrape = (id: number) => api.post<Article[]>(`/api/v1/sources/${id}/test`).then(r => r.data)

// Scrape
export const scrapeAll = () => api.post<Article[]>('/api/v1/scrape/all').then(r => r.data)

// Articles
export const getArticles = (sourceId?: number) => {
  const params = sourceId ? { source_id: sourceId } : {}
  return api.get<Article[]>('/api/v1/articles', { params }).then(r => r.data)
}

// Health
export const healthCheck = () => api.get('/api/v1/health').then(r => r.data)