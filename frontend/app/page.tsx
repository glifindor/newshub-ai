export default function Home() {
  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="container mx-auto px-4 py-16">
        {/* Header */}
        <header className="text-center mb-16">
          <h1 className="text-6xl font-bold text-gray-900 mb-4">
            🤖 NewsHub AI
          </h1>
          <p className="text-xl text-gray-600">
            Умные новости с AI-анализом для криптовалют и политики
          </p>
        </header>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-16">
          <div className="bg-white rounded-lg shadow-lg p-6 text-center">
            <div className="text-4xl font-bold text-blue-600 mb-2">1,234</div>
            <div className="text-gray-600">Новостей собрано</div>
          </div>
          <div className="bg-white rounded-lg shadow-lg p-6 text-center">
            <div className="text-4xl font-bold text-green-600 mb-2">567</div>
            <div className="text-gray-600">Опубликовано в Telegram</div>
          </div>
          <div className="bg-white rounded-lg shadow-lg p-6 text-center">
            <div className="text-4xl font-bold text-purple-600 mb-2">15</div>
            <div className="text-gray-600">Активных источников</div>
          </div>
        </div>

        {/* Channels */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-16">
          <div className="bg-white rounded-lg shadow-lg p-8">
            <div className="text-3xl mb-4">🔐</div>
            <h2 className="text-2xl font-bold mb-2">Crypto AI News</h2>
            <p className="text-gray-600 mb-4">
              Криптовалюты, IT и AI-анализ для инвесторов
            </p>
            <a
              href="https://t.me/crypto_ainews"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-block bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition"
            >
              Подписаться →
            </a>
          </div>

          <div className="bg-white rounded-lg shadow-lg p-8">
            <div className="text-3xl mb-4">🏛️</div>
            <h2 className="text-2xl font-bold mb-2">Kremlin Digest</h2>
            <p className="text-gray-600 mb-4">
              Политика России и мира с AI-разбором
            </p>
            <a
              href="https://t.me/kremlin_digest"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-block bg-red-600 text-white px-6 py-2 rounded-lg hover:bg-red-700 transition"
            >
              Подписаться →
            </a>
          </div>
        </div>

        {/* Latest News */}
        <div className="bg-white rounded-lg shadow-lg p-8">
          <h2 className="text-3xl font-bold mb-6">📰 Последние новости</h2>
          
          <div className="space-y-4">
            {/* Sample News Item */}
            <div className="border-l-4 border-blue-500 pl-4 py-2">
              <h3 className="font-bold text-lg mb-1">
                Bitcoin достиг нового ATH в $100,000
              </h3>
              <p className="text-gray-600 text-sm mb-2">
                Криптовалюта Bitcoin установила новый исторический максимум...
              </p>
              <div className="flex items-center gap-4 text-sm text-gray-500">
                <span>🔐 Crypto</span>
                <span>⭐ 9/10</span>
                <span>🕐 2 часа назад</span>
              </div>
            </div>

            <div className="border-l-4 border-red-500 pl-4 py-2">
              <h3 className="font-bold text-lg mb-1">
                Новое заявление Кремля по внешней политике
              </h3>
              <p className="text-gray-600 text-sm mb-2">
                Официальный представитель МИД России выступил с комментарием...
              </p>
              <div className="flex items-center gap-4 text-sm text-gray-500">
                <span>🏛️ Politics</span>
                <span>⭐ 8/10</span>
                <span>🕐 3 часа назад</span>
              </div>
            </div>
          </div>

          <div className="text-center mt-6">
            <a
              href="/news"
              className="inline-block bg-gray-800 text-white px-8 py-3 rounded-lg hover:bg-gray-900 transition"
            >
              Все новости →
            </a>
          </div>
        </div>

        {/* Footer */}
        <footer className="text-center mt-16 text-gray-600">
          <p>
            Powered by OpenRouter AI • Deployed on Docker
          </p>
          <p className="mt-2">
            <a href="/api/docs" className="text-blue-600 hover:underline">
              API Docs
            </a>
            {' • '}
            <a href="https://github.com/glifindor/newsportal" className="text-blue-600 hover:underline">
              GitHub
            </a>
          </p>
        </footer>
      </div>
    </main>
  )
}
