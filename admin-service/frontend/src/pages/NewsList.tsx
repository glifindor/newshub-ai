import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { newsAPI, News, PaginatedResponse } from '../api/api'

export default function NewsList() {
  const [news, setNews] = useState<PaginatedResponse<News> | null>(null)
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [status, setStatus] = useState<string>('')

  useEffect(() => {
    loadNews()
  }, [page, search, status])

  const loadNews = async () => {
    setLoading(true)
    try {
      const data = await newsAPI.getNews({
        page,
        page_size: 20,
        search: search || undefined,
        status: status || undefined,
      })
      setNews(data)
    } catch (error) {
      console.error('Failed to load news:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Удалить новость?')) return

    try {
      await newsAPI.deleteNews(id)
      loadNews()
    } catch (error) {
      console.error('Failed to delete news:', error)
      alert('Ошибка при удалении новости')
    }
  }

  const getStatusBadge = (status: string) => {
    const classes = {
      published: 'bg-green-100 text-green-800',
      draft: 'bg-yellow-100 text-yellow-800',
      archived: 'bg-gray-100 text-gray-800',
    }
    const labels = {
      published: 'Опубликовано',
      draft: 'Черновик',
      archived: 'Архив',
    }
    
    return (
      <span className={`px-2 py-1 text-xs rounded-full ${classes[status as keyof typeof classes]}`}>
        {labels[status as keyof typeof labels]}
      </span>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Новости</h1>
        <Link
          to="/news/new"
          className="bg-primary-600 hover:bg-primary-700 text-white px-6 py-2 rounded-md transition"
        >
          + Создать новость
        </Link>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-lg shadow p-4 flex gap-4">
        <input
          type="text"
          placeholder="Поиск..."
          value={search}
          onChange={(e) => {
            setSearch(e.target.value)
            setPage(1)
          }}
          className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
        />
        <select
          value={status}
          onChange={(e) => {
            setStatus(e.target.value)
            setPage(1)
          }}
          className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
        >
          <option value="">Все статусы</option>
          <option value="published">Опубликовано</option>
          <option value="draft">Черновики</option>
          <option value="archived">Архив</option>
        </select>
      </div>

      {/* News List */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        {loading ? (
          <div className="text-center py-12">Загрузка...</div>
        ) : news && news.data.length > 0 ? (
          <>
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Заголовок
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Автор
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Статус
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Просмотры
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Дата
                  </th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Действия
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {news.data.map((item) => (
                  <tr key={item.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4">
                      <div className="text-sm font-medium text-gray-900">{item.title}</div>
                      {item.category && (
                        <div className="text-sm text-gray-500">{item.category}</div>
                      )}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500">{item.author_name}</td>
                    <td className="px-6 py-4">{getStatusBadge(item.status)}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{item.view_count}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">
                      {new Date(item.created_at).toLocaleDateString('ru-RU')}
                    </td>
                    <td className="px-6 py-4 text-right text-sm font-medium space-x-2">
                      <Link
                        to={`/news/${item.id}`}
                        className="text-primary-600 hover:text-primary-900"
                      >
                        Редактировать
                      </Link>
                      <button
                        onClick={() => handleDelete(item.id)}
                        className="text-red-600 hover:text-red-900"
                      >
                        Удалить
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>

            {/* Pagination */}
            <div className="bg-gray-50 px-6 py-3 flex items-center justify-between border-t">
              <div className="text-sm text-gray-700">
                Показано {news.data.length} из {news.total} новостей
              </div>
              <div className="flex space-x-2">
                <button
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  disabled={page === 1}
                  className="px-3 py-1 border rounded-md disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100"
                >
                  Назад
                </button>
                <span className="px-3 py-1">
                  Страница {page} из {news.total_pages}
                </span>
                <button
                  onClick={() => setPage((p) => p + 1)}
                  disabled={page >= news.total_pages}
                  className="px-3 py-1 border rounded-md disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100"
                >
                  Вперед
                </button>
              </div>
            </div>
          </>
        ) : (
          <div className="text-center py-12 text-gray-500">Новости не найдены</div>
        )}
      </div>
    </div>
  )
}
