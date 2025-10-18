/**
 * Type definitions for NewsHub AI Frontend
 */

export type NewsChannel = 'crypto' | 'politics';

export type NewsStatus = 'pending' | 'analyzed' | 'published' | 'rejected';

export interface NewsItem {
  id: string;
  title: string;
  content: string;
  url: string;
  source_name: string;
  category: NewsChannel;
  status: NewsStatus;
  ai_summary?: string;
  ai_insights?: string[] | string;
  ai_hashtags?: string[];
  image_url?: string;
  relevance_score?: number;
  requires_moderation: boolean;
  telegram_message_id?: number;
  telegram_channel?: string;
  published_at?: string;
  created_at: string;
  updated_at?: string;
}

export interface NewsSource {
  id: string;
  name: string;
  type: 'rss' | 'api' | 'scraping';
  url: string;
  category: NewsChannel;
  keywords: string[];
  is_active: boolean;
  last_fetch_at?: string;
  created_at: string;
}

export interface User {
  id: string;
  email: string;
  username: string;
  role: 'admin' | 'moderator';
  is_active: boolean;
  created_at: string;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  per_page: number;
  pages: number;
}

export interface ApiError {
  detail: string;
  status?: number;
}

export interface LoginCredentials {
  username: string;
  password: string;
}

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  token_type: string;
}

export interface NewsFilters {
  search?: string;
  category?: NewsChannel;
  status?: NewsStatus;
  page?: number;
  per_page?: number;
  sort_by?: 'created_at' | 'relevance_score' | 'published_at';
  order?: 'asc' | 'desc';
}

export interface SourceFormData {
  name: string;
  type: 'rss' | 'api' | 'scraping';
  url: string;
  category: NewsChannel;
  keywords: string[];
  is_active: boolean;
}

export interface PipelineResult {
  total: number;
  posted?: number;
  analyzed?: number;
  collected?: number;
  failed?: number;
}

export interface WebSocketMessage {
  type: 'news_published' | 'news_analyzed' | 'news_collected' | 'error';
  data: NewsItem | PipelineResult | { message: string };
  timestamp: string;
}
