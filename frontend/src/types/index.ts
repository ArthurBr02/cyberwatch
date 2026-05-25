export interface Source {
  id: number
  name: string
  url: string
  fetch_type: 'html' | 'json'
  llm_rules: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface Article {
  id: number
  source_id: number
  title: string
  url: string
  summary: string
  published_at: string | null
  is_sent: boolean
  created_at: string
  updated_at: string
}

export interface CreateSourcePayload {
  name: string
  url: string
  fetch_type: 'html' | 'json'
}

export interface UpdateSourcePayload {
  name: string
  url: string
  fetch_type: 'html' | 'json'
  is_active: boolean
}