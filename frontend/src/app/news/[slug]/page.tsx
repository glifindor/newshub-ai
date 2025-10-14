import { notFound } from 'next/navigation'
import { getNewsBySlug, incrementViewCount } from '@/lib/api'
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'
import type { Metadata } from 'next'

interface NewsPageProps {
  params: { slug: string }
}

export async function generateMetadata({ params }: NewsPageProps): Promise<Metadata> {
  const news = await getNewsBySlug(params.slug)
  
  if (!news) {
    return { title: 'Новость не найдена' }
  }

  return {
    title: `${news.title} - Новостной портал`,
    description: news.summary || news.content.substring(0, 160),
  }
}

export default async function NewsPage({ params }: NewsPageProps) {
  const news = await getNewsBySlug(params.slug)
  
  if (!news) {
    notFound()
  }

  // Increment view count (non-blocking)
  incrementViewCount(news.id).catch(console.error)

  return (
    <article className="bg-white">
      <div className="container mx-auto px-4 py-8 max-w-4xl">
        {/* Header */}
        <header className="mb-8">
          {news.category && (
            <span className="inline-block bg-primary-100 text-primary-700 text-sm font-semibold px-4 py-1 rounded-full mb-4">
              {news.category}
            </span>
          )}
          
          <h1 className="text-4xl font-bold mb-4 text-gray-900">{news.title}</h1>
          
          <div className="flex items-center gap-4 text-gray-600 text-sm mb-4">
            <span className="flex items-center gap-1">
              👤 {news.author_name}
            </span>
            <span className="flex items-center gap-1">
              📅 {format(new Date(news.created_at), 'dd MMMM yyyy, HH:mm', { locale: ru })}
            </span>
            <span className="flex items-center gap-1">
              👁️ {news.view_count} просмотров
            </span>
          </div>

          {news.tags && news.tags.length > 0 && (
            <div className="flex flex-wrap gap-2">
              {news.tags.map((tag) => (
                <span key={tag} className="text-sm text-primary-600">#{tag}</span>
              ))}
            </div>
          )}
        </header>

        {/* Cover Image */}
        {news.cover_image && (
          <div className="mb-8">
            <img
              src={news.cover_image}
              alt={news.title}
              className="w-full rounded-lg shadow-lg"
            />
          </div>
        )}

        {/* Summary */}
        {news.summary && (
          <div className="text-xl text-gray-700 mb-8 p-6 bg-gray-50 rounded-lg border-l-4 border-primary-500">
            {news.summary}
          </div>
        )}

        {/* Content */}
        <div 
          className="news-content prose prose-lg max-w-none"
          dangerouslySetInnerHTML={{ __html: news.content }}
        />

        {/* Footer */}
        <footer className="mt-12 pt-8 border-t">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-gray-600">Автор: <strong>{news.author_name}</strong></p>
              <p className="text-sm text-gray-500">
                Опубликовано: {format(new Date(news.created_at), 'dd MMMM yyyy', { locale: ru })}
              </p>
            </div>
            <div className="flex gap-2">
              <button className="px-4 py-2 bg-primary-100 text-primary-700 rounded-lg hover:bg-primary-200 transition">
                Поделиться
              </button>
            </div>
          </div>
        </footer>
      </div>
    </article>
  )
}
