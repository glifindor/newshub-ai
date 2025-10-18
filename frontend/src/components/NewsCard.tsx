/**
 * NewsCard Component - карточка новости
 */
import React from 'react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { motion } from 'framer-motion';
import { FiExternalLink, FiClock, FiTrendingUp } from 'react-icons/fi';
import { SiTelegram } from 'react-icons/si';
import type { NewsItem } from '@/types';
import { cn } from '@/utils/cn';

interface NewsCardProps {
  news: NewsItem;
  onClick?: () => void;
  showActions?: boolean;
  onApprove?: (id: string) => void;
  onReject?: (id: string) => void;
}

export const NewsCard: React.FC<NewsCardProps> = ({
  news,
  onClick,
  showActions = false,
  onApprove,
  onReject,
}) => {
  const categoryConfig = {
    crypto: {
      color: 'text-crypto-600 dark:text-crypto-400',
      bg: 'bg-crypto-100 dark:bg-crypto-900/20',
      icon: '🔐',
    },
    politics: {
      color: 'text-politics-600 dark:text-politics-400',
      bg: 'bg-politics-100 dark:bg-politics-900/20',
      icon: '🏛️',
    },
  };

  const statusConfig = {
    pending: { label: 'Ожидает', color: 'bg-gray-500' },
    analyzed: { label: 'Проанализирована', color: 'bg-blue-500' },
    published: { label: 'Опубликована', color: 'bg-green-500' },
    rejected: { label: 'Отклонена', color: 'bg-red-500' },
  };

  const config = categoryConfig[news.category];
  const status = statusConfig[news.status];

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      className={cn(
        'group relative overflow-hidden rounded-xl border border-gray-200 dark:border-gray-800',
        'bg-white dark:bg-gray-900 shadow-sm hover:shadow-md transition-all duration-200',
        onClick && 'cursor-pointer'
      )}
      onClick={onClick}
    >
      {/* Image */}
      {news.image_url && (
        <div className="relative h-48 w-full overflow-hidden bg-gray-100 dark:bg-gray-800">
          <img
            src={news.image_url}
            alt={news.title}
            className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
            loading="lazy"
          />
          {/* Category badge */}
          <div className={cn('absolute top-3 left-3 px-3 py-1 rounded-full text-sm font-medium', config.bg)}>
            <span className="mr-1">{config.icon}</span>
            <span className={config.color}>{news.category}</span>
          </div>
        </div>
      )}

      <div className="p-5">
        {/* Status & Relevance */}
        <div className="mb-3 flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <span className={cn('h-2 w-2 rounded-full', status.color)} />
            <span className="text-xs text-gray-600 dark:text-gray-400">{status.label}</span>
          </div>
          {news.relevance_score && (
            <div className="flex items-center space-x-1 text-xs text-gray-600 dark:text-gray-400">
              <FiTrendingUp className="h-3 w-3" />
              <span>{news.relevance_score.toFixed(1)}/10</span>
            </div>
          )}
        </div>

        {/* Title */}
        <h3 className="mb-2 text-lg font-bold text-gray-900 dark:text-white line-clamp-2 group-hover:text-primary-600 dark:group-hover:text-primary-400 transition-colors">
          {news.title}
        </h3>

        {/* AI Summary */}
        {news.ai_summary && (
          <p className="mb-3 text-sm text-gray-600 dark:text-gray-400 line-clamp-3">
            {news.ai_summary}
          </p>
        )}

        {/* AI Insights */}
        {news.ai_insights && (
          <div className="mb-3 space-y-1">
            {(Array.isArray(news.ai_insights) ? news.ai_insights : [news.ai_insights])
              .slice(0, 2)
              .map((insight, idx) => (
                <div key={idx} className="flex items-start space-x-2 text-xs text-gray-600 dark:text-gray-400">
                  <span className="mt-0.5">•</span>
                  <span className="line-clamp-1">{insight}</span>
                </div>
              ))}
          </div>
        )}

        {/* Hashtags */}
        {news.ai_hashtags && news.ai_hashtags.length > 0 && (
          <div className="mb-3 flex flex-wrap gap-1">
            {news.ai_hashtags.slice(0, 3).map((tag, idx) => (
              <span
                key={idx}
                className="inline-block px-2 py-1 text-xs rounded bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300"
              >
                {tag}
              </span>
            ))}
          </div>
        )}

        {/* Footer */}
        <div className="flex items-center justify-between pt-3 border-t border-gray-100 dark:border-gray-800">
          <div className="flex items-center space-x-2 text-xs text-gray-500 dark:text-gray-500">
            <FiClock className="h-3 w-3" />
            <time>{format(new Date(news.created_at), 'dd MMM yyyy, HH:mm', { locale: ru })}</time>
          </div>

          <div className="flex items-center space-x-2">
            {news.telegram_message_id && (
              <a
                href={`https://t.me/${news.telegram_channel?.replace('@', '')}/${news.telegram_message_id}`}
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary-600 hover:text-primary-700 dark:text-primary-400"
                onClick={(e) => e.stopPropagation()}
              >
                <SiTelegram className="h-4 w-4" />
              </a>
            )}
            <a
              href={news.url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
              onClick={(e) => e.stopPropagation()}
            >
              <FiExternalLink className="h-4 w-4" />
            </a>
          </div>
        </div>

        {/* Action buttons for admin */}
        {showActions && (
          <div className="mt-4 flex items-center space-x-2">
            {news.status === 'analyzed' && (
              <>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    onApprove?.(news.id);
                  }}
                  className="flex-1 px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-lg hover:bg-green-700 transition-colors"
                >
                  ✅ Одобрить
                </button>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    onReject?.(news.id);
                  }}
                  className="flex-1 px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors"
                >
                  ❌ Отклонить
                </button>
              </>
            )}
          </div>
        )}
      </div>
    </motion.div>
  );
};

export default NewsCard;
