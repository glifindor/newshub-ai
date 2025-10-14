import axios from 'axios'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export interface News {
  id: string
  title: string
  slug: string
  summary?: string
  content: string
  cover_image?: string
  category?: string
  tags?: string[]
  author_id: string
  author_name: string
  status: string
  view_count: number
  created_at: string
  updated_at: string
  published_at?: string
}

export interface PaginatedNews {
  data: News[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface GetNewsParams {
  page?: number
  page_size?: number
  status?: string
  category?: string
  tag?: string
  search?: string
}

export async function getNews(params: GetNewsParams = {}): Promise<PaginatedNews> {
  try {
    const response = await api.get('/api/v1/news', { params })
    return response.data
  } catch (error) {
    console.error('Failed to fetch news:', error)
    return {
      data: [],
      total: 0,
      page: 1,
      page_size: 20,
      total_pages: 0,
    }
  }
}

export async function getNewsBySlug(slug: string): Promise<News | null> {
  try {
    const response = await api.get(`/api/v1/news/slug/${slug}`)
    return response.data
  } catch (error) {
    console.error('Failed to fetch news by slug:', error)
    return null
  }
}

export async function incrementViewCount(id: string): Promise<void> {
  try {
    await api.post(`/api/v1/news/${id}/view`)
  } catch (error) {
    console.error('Failed to increment view count:', error)
  }
}

export async function getCategories(): Promise<string[]> {
  try {
    const response = await api.get('/api/v1/categories')
    return response.data
  } catch (error) {
    console.error('Failed to fetch categories:', error)
    return []
  }
}
