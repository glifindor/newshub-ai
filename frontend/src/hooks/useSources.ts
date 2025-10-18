/**
 * Custom hooks for news sources
 */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { sourcesApi } from '@/lib/api';
import type { SourceFormData } from '@/types';
import { toast } from 'react-hot-toast';

export const sourcesKeys = {
  all: ['sources'] as const,
  lists: () => [...sourcesKeys.all, 'list'] as const,
  list: () => [...sourcesKeys.lists()] as const,
  details: () => [...sourcesKeys.all, 'detail'] as const,
  detail: (id: string) => [...sourcesKeys.details(), id] as const,
};

/**
 * Fetch all sources
 */
export function useSources() {
  return useQuery({
    queryKey: sourcesKeys.list(),
    queryFn: () => sourcesApi.getAll(),
    staleTime: 60000, // 1 minute
  });
}

/**
 * Fetch single source
 */
export function useSource(id: string, enabled = true) {
  return useQuery({
    queryKey: sourcesKeys.detail(id),
    queryFn: () => sourcesApi.getById(id),
    enabled: enabled && !!id,
  });
}

/**
 * Create source mutation
 */
export function useCreateSource() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: SourceFormData) => sourcesApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: sourcesKeys.lists() });
      toast.success('✅ Источник создан успешно!');
    },
    onError: (error: any) => {
      toast.error(`❌ Ошибка: ${error.response?.data?.detail || error.message}`);
    },
  });
}

/**
 * Update source mutation
 */
export function useUpdateSource() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<SourceFormData> }) =>
      sourcesApi.update(id, data),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: sourcesKeys.lists() });
      queryClient.setQueryData(sourcesKeys.detail(data.id), data);
      toast.success('✅ Источник обновлён');
    },
    onError: (error: any) => {
      toast.error(`❌ Ошибка: ${error.response?.data?.detail || error.message}`);
    },
  });
}

/**
 * Delete source mutation
 */
export function useDeleteSource() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => sourcesApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: sourcesKeys.lists() });
      toast.success('🗑️ Источник удалён');
    },
    onError: (error: any) => {
      toast.error(`❌ Ошибка: ${error.response?.data?.detail || error.message}`);
    },
  });
}

/**
 * Toggle source active status
 */
export function useToggleSource() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => sourcesApi.toggle(id),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: sourcesKeys.lists() });
      queryClient.setQueryData(sourcesKeys.detail(data.id), data);
      toast.success(data.is_active ? '✅ Источник активирован' : '⏸️ Источник деактивирован');
    },
    onError: (error: any) => {
      toast.error(`❌ Ошибка: ${error.response?.data?.detail || error.message}`);
    },
  });
}
