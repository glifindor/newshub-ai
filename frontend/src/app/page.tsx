import NewsGrid from '@/components/NewsGrid'
import CategoryNav from '@/components/CategoryNav'
import { getNews } from '@/lib/api'

export const revalidate = 60 // Revalidate every 60 seconds

export default async function HomePage() {
  const news = await getNews({ page: 1, page_size: 12, status: 'published' })

  return (
    <div className="bg-gray-50">
      {/* Hero Section */}
      <section className="bg-gradient-to-r from-primary-600 to-primary-800 text-white py-16">
        <div className="container mx-auto px-4">
          <h1 className="text-5xl font-bold mb-4">Новостной портал</h1>
          <p className="text-xl text-primary-100">Будьте в курсе самых важных событий</p>
        </div>
      </section>

      {/* Category Navigation */}
      <CategoryNav />

      {/* Latest News */}
      <section className="container mx-auto px-4 py-12">
        <h2 className="text-3xl font-bold mb-8">Последние новости</h2>
        <NewsGrid news={news.data} />
        
        {news.total > 12 && (
          <div className="text-center mt-8">
            <a
              href="/news"
              className="inline-block bg-primary-600 hover:bg-primary-700 text-white px-8 py-3 rounded-lg font-semibold transition"
            >
              Все новости
            </a>
          </div>
        )}
      </section>
    </div>
  )
}
