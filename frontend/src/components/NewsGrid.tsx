import Link from 'next/link'
import { News } from '@/lib/api'
import { formatDistanceToNow } from 'date-fns'
import { ru } from 'date-fns/locale'

interface NewsGridProps {
  news: News[]
}

export default function NewsGrid({ news }: NewsGridProps) {
  if (news.length === 0) {
    return (
      <div className="text-center py-12 text-gray-500">
        <p className="text-xl">Новостей пока нет</p>
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {news.map((item) => (
        <article key={item.id} className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-xl transition">
          {item.cover_image && (
            <div className="h-48 bg-gray-200 overflow-hidden">
              <img
                src={item.cover_image}
                alt={item.title}
                className="w-full h-full object-cover hover:scale-105 transition-transform duration-300"
              />
            </div>
          )}
          
          <div className="p-6">
            {item.category && (
              <span className="inline-block bg-primary-100 text-primary-700 text-xs font-semibold px-3 py-1 rounded-full mb-3">
                {item.category}
              </span>
            )}
            
            <h3 className="text-xl font-bold mb-2 line-clamp-2 hover:text-primary-600 transition">
              <Link href={`/news/${item.slug}`}>
                {item.title}
              </Link>
            </h3>
            
            {item.summary && (
              <p className="text-gray-600 mb-4 line-clamp-3">{item.summary}</p>
            )}
            
            <div className="flex items-center justify-between text-sm text-gray-500">
              <span>{item.author_name}</span>
              <time>{formatDistanceToNow(new Date(item.created_at), { addSuffix: true, locale: ru })}</time>
            </div>
            
            <div className="flex items-center gap-4 mt-3 text-sm text-gray-500">
              <span>👁️ {item.view_count}</span>
              {item.tags && item.tags.length > 0 && (
                <div className="flex gap-1">
                  {item.tags.slice(0, 2).map((tag) => (
                    <span key={tag} className="text-xs text-primary-600">#{tag}</span>
                  ))}
                </div>
              )}
            </div>
          </div>
        </article>
      ))}
    </div>
  )
}
