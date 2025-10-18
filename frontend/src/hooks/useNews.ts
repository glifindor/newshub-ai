/**
 * Custom hooks for news operations
 */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { newsApi } from '@/lib/api';
import type { NewsFilters, NewsItem } from '@/types';
import { toast } from 'react-hot-toast';

// ============================================================================
// Query Keys
// ============================================================================

export const newsKeys = {
  all: ['news'] as const,
  lists: () => [...newsKeys.all, 'list'] as const,
  list: (filters?: NewsFilters) => [...newsKeys.lists(), filters] as const,
  details: () => [...newsKeys.all, 'detail'] as const,
  detail: (id: string) => [...newsKeys.details(), id] as const,
  public: (filters?: NewsFilters) => ['news', 'public', filters] as const,
  channel: (channel: string, filters?: NewsFilters) => 
    ['news', 'channel', channel, filters] as const,
};

// ============================================================================
// Hooks
// ============================================================================

/**
 * Fetch all news (admin)
 */
export function useNews(filters?: NewsFilters) {
  return useQuery({
    queryKey: newsKeys.list(filters),
    queryFn: () => newsApi.getAll(filters),
    staleTime: 30000, // 30 seconds
  });
}

/**
 * Fetch single news item
 */
export function useNewsItem(id: string, enabled = true) {
  return useQuery({
    queryKey: newsKeys.detail(id),
    queryFn: () => newsApi.getById(id),
    enabled: enabled && !!id,
  });
}

/**
 * Fetch public news
 */
export function usePublicNews(filters?: NewsFilters) {
  return useQuery({
    queryKey: newsKeys.public(filters),
    queryFn: () => newsApi.getPublic(filters),
    staleTime: 60000, // 1 minute
  });
}

/**
 * Fetch news by channel
 */
export function useChannelNews(channel: string, filters?: NewsFilters) {
  return useQuery({
    queryKey: newsKeys.channel(channel, filters),
    queryFn: () => newsApi.getByChannel(channel, filters),
    staleTime: 60000,
  });
}

/**
 * Approve news mutation
 */
export function useApproveNews() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => newsApi.approve(id),
    onSuccess: (data) => {
      // Invalidate and refetch
      queryClient.invalidateQueries({ queryKey: newsKeys.lists() });
      queryClient.setQueryData(newsKeys.detail(data.id), data);
      toast.success('✅ Новость одобрена и опубликована!');
    },
    onError: (error: any) => {
      toast.error(`❌ Ошибка: ${error.response?.data?.detail || error.message}`);
    },
  });
}

/**
 * Reject news mutation
 */
export function useRejectNews() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => newsApi.reject(id),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: newsKeys.lists() });
      queryClient.setQueryData(newsKeys.detail(data.id), data);
      toast.success('❌ Новость отклонена');
    },
    onError: (error: any) => {
      toast.error(`❌ Ошибка: ${error.response?.data?.detail || error.message}`);
    },
  });
}

/**
 * Delete news mutation
 */
export function useDeleteNews() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => newsApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: newsKeys.lists() });
      toast.success('🗑️ Новость удалена');
    },
    onError: (error: any) => {
      toast.error(`❌ Ошибка: ${error.response?.data?.detail || error.message}`);
    },
  });
}
