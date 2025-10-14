import axios, { AxiosError } from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/admin'

export interface User {
  id: string
  username: string
  email: string
  role: string
  created_at: string
}

export interface News {
  id: string
  title: string
  slug: string
  summary: string
  content: string
  author_id: string
  author_name: string
  category_id?: string
  category?: string
  tags: string[]
  image_url: string
  status: 'draft' | 'published' | 'archived'
  view_count: number
  published_at?: string
  created_at: string
  updated_at: string
}

export interface SEOMeta {
  id: string
  news_id: string
  slug: string
  title: string
  description: string
  keywords: string
  canonical_url: string
  og_title: string
  og_description: string
  og_image: string
  robots_index: boolean
  robots_follow: boolean
  sitemap_change_freq: string
  sitemap_priority: number
  schema_data?: any
  created_at: string
  updated_at: string
}

export interface DashboardStats {
  total_news: number
  published_news: number
  draft_news: number
  total_views: number
  today_views: number
  popular_news: News[]
  recent_news: News[]
  category_stats: Record<string, number>
  views_trend: Array<{ date: string; views: number }>
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface CreateNewsRequest {
  title: string
  slug: string
  summary: string
  content: string
  category_id?: string
  tags: string[]
  image_url: string
  status: 'draft' | 'published'
}

export interface UpdateNewsRequest {
  title?: string
  summary?: string
  content?: string
  category_id?: string
  tags?: string[]
  image_url?: string
  status?: 'draft' | 'published' | 'archived'
}

export interface UpdateSEORequest {
  title?: string
  description?: string
  keywords?: string
  og_title?: string
  og_description?: string
  og_image?: string
  robots_index?: boolean
  robots_follow?: boolean
  sitemap_change_freq?: string
  sitemap_priority?: number
}

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor для добавления токена
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor для обработки ошибок
api.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Auth API
export const authAPI = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/login', data)
    return response.data
  },

  getCurrentUser: async (): Promise<User> => {
    const response = await api.get<User>('/me')
    return response.data
  },
}

// News API
export const newsAPI = {
  getNews: async (params: {
    page?: number
    page_size?: number
    search?: string
    status?: string
    category?: string
    tag?: string
  }): Promise<PaginatedResponse<News>> => {
    const response = await api.get<PaginatedResponse<News>>('/news', { params })
    return response.data
  },

  getNewsById: async (id: string): Promise<News> => {
    const response = await api.get<News>(`/news/${id}`)
    return response.data
  },

  createNews: async (data: CreateNewsRequest): Promise<News> => {
    const response = await api.post<News>('/news', data)
    return response.data
  },

  updateNews: async (id: string, data: UpdateNewsRequest): Promise<News> => {
    const response = await api.put<News>(`/news/${id}`, data)
    return response.data
  },

  deleteNews: async (id: string): Promise<void> => {
    await api.delete(`/news/${id}`)
  },
}

// SEO API
export const seoAPI = {
  getNewsSEO: async (newsId: string): Promise<SEOMeta | null> => {
    try {
      const response = await api.get<SEOMeta>(`/news/${newsId}/seo`)
      return response.data
    } catch (error) {
      if ((error as AxiosError).response?.status === 404) {
        return null
      }
      throw error
    }
  },

  updateNewsSEO: async (newsId: string, data: UpdateSEORequest): Promise<SEOMeta> => {
    const response = await api.put<SEOMeta>(`/news/${newsId}/seo`, data)
    return response.data
  },
}

// Dashboard API
export const dashboardAPI = {
  getStats: async (): Promise<DashboardStats> => {
    const response = await api.get<DashboardStats>('/dashboard')
    return response.data
  },
}

export default api
