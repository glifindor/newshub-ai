import { useEffect, useState } from 'react'
import { dashboardAPI, DashboardStats } from '../api/api'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'

export default function Dashboard() {
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadStats()
  }, [])

  const loadStats = async () => {
    try {
      const data = await dashboardAPI.getStats()
      setStats(data)
    } catch (error) {
      console.error('Failed to load stats:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="text-center py-12">Загрузка...</div>
  }

  if (!stats) {
    return <div className="text-center py-12 text-red-600">Ошибка загрузки статистики</div>
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">Дашборд</h1>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatsCard title="Всего новостей" value={stats.total_news} icon="📰" />
        <StatsCard title="Опубликовано" value={stats.published_news} icon="✅" color="green" />
        <StatsCard title="Черновики" value={stats.draft_news} icon="📝" color="yellow" />
        <StatsCard title="Всего просмотров" value={stats.total_views} icon="👁️" color="blue" />
      </div>

      {/* Views Trend Chart */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-xl font-bold mb-4">Динамика просмотров</h2>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={stats.views_trend}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" />
            <YAxis />
            <Tooltip />
            <Legend />
            <Line type="monotone" dataKey="views" stroke="#0ea5e9" strokeWidth={2} />
          </LineChart>
        </ResponsiveContainer>
      </div>

      {/* Recent and Popular News */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Popular News */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-bold mb-4">Популярные новости</h2>
          <div className="space-y-3">
            {stats.popular_news.slice(0, 5).map((news) => (
              <div key={news.id} className="flex justify-between items-center border-b pb-2">
                <div className="flex-1">
                  <p className="font-medium truncate">{news.title}</p>
                  <p className="text-sm text-gray-500">{news.view_count} просмотров</p>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Recent News */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-bold mb-4">Последние новости</h2>
          <div className="space-y-3">
            {stats.recent_news.slice(0, 5).map((news) => (
              <div key={news.id} className="border-b pb-2">
                <p className="font-medium truncate">{news.title}</p>
                <p className="text-sm text-gray-500">
                  {new Date(news.created_at).toLocaleDateString('ru-RU')}
                </p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

function StatsCard({ 
  title, 
  value, 
  icon, 
  color = 'primary' 
}: { 
  title: string
  value: number
  icon: string
  color?: string
}) {
  const colorClasses = {
    primary: 'bg-primary-50 border-primary-200',
    green: 'bg-green-50 border-green-200',
    yellow: 'bg-yellow-50 border-yellow-200',
    blue: 'bg-blue-50 border-blue-200',
  }

  return (
    <div className={`${colorClasses[color as keyof typeof colorClasses] || colorClasses.primary} border rounded-lg p-6`}>
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm text-gray-600">{title}</p>
          <p className="text-3xl font-bold mt-2">{value.toLocaleString()}</p>
        </div>
        <div className="text-4xl">{icon}</div>
      </div>
    </div>
  )
}
