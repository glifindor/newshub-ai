/**
 * API Client for NewsHub AI Backend
 */
import axios, { AxiosError, AxiosInstance } from 'axios';
import type {
  NewsItem,
  NewsSource,
  PaginatedResponse,
  NewsFilters,
  SourceFormData,
  LoginCredentials,
  AuthTokens,
  PipelineResult,
} from '@/types';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000/api/v1';

// Create axios instance
export const apiClient: AxiosInstance = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000,
});

// Request interceptor for adding auth token
apiClient.interceptors.request.use(
  (config) => {
    // Get token from localStorage (client-side only)
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('access_token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for handling errors
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    if (error.response?.status === 401) {
      // Token expired, try to refresh
      if (typeof window !== 'undefined') {
        const refreshToken = localStorage.getItem('refresh_token');
        if (refreshToken) {
          try {
            const { data } = await axios.post(`${API_URL}/auth/refresh`, {
              refresh_token: refreshToken,
            });
            localStorage.setItem('access_token', data.access_token);
            // Retry original request
            if (error.config) {
              error.config.headers.Authorization = `Bearer ${data.access_token}`;
              return axios(error.config);
            }
          } catch (refreshError) {
            // Refresh failed, logout
            localStorage.removeItem('access_token');
            localStorage.removeItem('refresh_token');
            window.location.href = '/admin/login';
          }
        }
      }
    }
    return Promise.reject(error);
  }
);

// ============================================================================
// AUTH API
// ============================================================================

export const authApi = {
  login: async (credentials: LoginCredentials): Promise<AuthTokens> => {
    const { data } = await apiClient.post<AuthTokens>('/auth/login', credentials);
    // Store tokens
    if (typeof window !== 'undefined') {
      localStorage.setItem('access_token', data.access_token);
      localStorage.setItem('refresh_token', data.refresh_token);
    }
    return data;
  },

  logout: async (): Promise<void> => {
    await apiClient.post('/auth/logout');
    if (typeof window !== 'undefined') {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
    }
  },

  refresh: async (refreshToken: string): Promise<AuthTokens> => {
    const { data } = await apiClient.post<AuthTokens>('/auth/refresh', {
      refresh_token: refreshToken,
    });
    return data;
  },
};

// ============================================================================
// NEWS API
// ============================================================================

export const newsApi = {
  getAll: async (filters?: NewsFilters): Promise<PaginatedResponse<NewsItem>> => {
    const { data } = await apiClient.get<PaginatedResponse<NewsItem>>('/news/', {
      params: filters,
    });
    return data;
  },

  getById: async (id: string): Promise<NewsItem> => {
    const { data } = await apiClient.get<NewsItem>(`/news/${id}`);
    return data;
  },

  approve: async (id: string): Promise<NewsItem> => {
    const { data } = await apiClient.post<NewsItem>(`/news/${id}/approve`);
    return data;
  },

  reject: async (id: string): Promise<NewsItem> => {
    const { data } = await apiClient.post<NewsItem>(`/news/${id}/reject`);
    return data;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/news/${id}`);
  },

  // Public endpoints
  getPublic: async (filters?: NewsFilters): Promise<PaginatedResponse<NewsItem>> => {
    const { data } = await apiClient.get<PaginatedResponse<NewsItem>>('/news/public', {
      params: filters,
    });
    return data;
  },

  getByChannel: async (
    channel: string,
    filters?: NewsFilters
  ): Promise<PaginatedResponse<NewsItem>> => {
    const { data } = await apiClient.get<PaginatedResponse<NewsItem>>(
      `/news/public/${channel}`,
      { params: filters }
    );
    return data;
  },
};

// ============================================================================
// SOURCES API
// ============================================================================

export const sourcesApi = {
  getAll: async (): Promise<NewsSource[]> => {
    const { data } = await apiClient.get<NewsSource[]>('/sources/');
    return data;
  },

  getById: async (id: string): Promise<NewsSource> => {
    const { data } = await apiClient.get<NewsSource>(`/sources/${id}`);
    return data;
  },

  create: async (source: SourceFormData): Promise<NewsSource> => {
    const { data } = await apiClient.post<NewsSource>('/sources/', source);
    return data;
  },

  update: async (id: string, source: Partial<SourceFormData>): Promise<NewsSource> => {
    const { data } = await apiClient.put<NewsSource>(`/sources/${id}`, source);
    return data;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/sources/${id}`);
  },

  toggle: async (id: string): Promise<NewsSource> => {
    const { data } = await apiClient.post<NewsSource>(`/sources/${id}/toggle`);
    return data;
  },
};

// ============================================================================
// PIPELINE API
// ============================================================================

export const pipelineApi = {
  collect: async (channel?: string): Promise<PipelineResult> => {
    const { data } = await apiClient.post<PipelineResult>('/pipeline/collect', null, {
      params: { channel },
    });
    return data;
  },

  analyze: async (limit?: number): Promise<PipelineResult> => {
    const { data } = await apiClient.post<PipelineResult>('/pipeline/analyze', null, {
      params: { limit },
    });
    return data;
  },

  post: async (limit?: number, channel?: string): Promise<PipelineResult> => {
    const { data } = await apiClient.post<PipelineResult>('/pipeline/post', null, {
      params: { limit, channel },
    });
    return data;
  },

  runFull: async (channel?: string): Promise<any> => {
    const { data } = await apiClient.post('/pipeline/pipeline', null, {
      params: { channel },
    });
    return data;
  },
};

// ============================================================================
// TELEGRAM API
// ============================================================================

export const telegramApi = {
  test: async (channel: string, title: string, content: string): Promise<any> => {
    const { data } = await apiClient.post('/telegram/test', {
      channel,
      title,
      content,
    });
    return data;
  },
};

export default apiClient;
