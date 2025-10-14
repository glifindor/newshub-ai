import NewsGrid from '@/components/NewsGrid'
import { getNews } from '@/lib/api'

interface NewsListPageProps {
  searchParams: { page?: string; category?: string; search?: string }
}

export default async function NewsListPage({ searchParams }: NewsListPageProps) {
  const page = parseInt(searchParams.page || '1')
  const category = searchParams.category
  const search = searchParams.search

  const news = await getNews({
    page,
    page_size: 20,
    status: 'published',
    category,
    search,
  })

  return (
    <div className="bg-gray-50 min-h-screen">
      <div className="container mx-auto px-4 py-12">
        <h1 className="text-4xl font-bold mb-8">
          {category ? `Категория: ${category}` : 'Все новости'}
        </h1>

        {search && (
          <p className="text-gray-600 mb-6">
            Результаты поиска: <strong>{search}</strong>
          </p>
        )}

        <NewsGrid news={news.data} />

        {/* Pagination */}
        {news.total_pages > 1 && (
          <div className="flex justify-center gap-2 mt-8">
            {page > 1 && (
              <a
                href={`/news?page=${page - 1}${category ? `&category=${category}` : ''}${search ? `&search=${search}` : ''}`}
                className="px-4 py-2 bg-white border rounded-lg hover:bg-gray-50"
              >
                ← Назад
              </a>
            )}
            
            <span className="px-4 py-2 bg-primary-600 text-white rounded-lg">
              Страница {page} из {news.total_pages}
            </span>
            
            {page < news.total_pages && (
              <a
                href={`/news?page=${page + 1}${category ? `&category=${category}` : ''}${search ? `&search=${search}` : ''}`}
                className="px-4 py-2 bg-white border rounded-lg hover:bg-gray-50"
              >
                Вперед →
              </a>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
