import Link from 'next/link'

export default function Header() {
  return (
    <header className="bg-white shadow-sm sticky top-0 z-50">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <Link href="/" className="text-2xl font-bold text-primary-600">
            📰 Новостной портал
          </Link>
          
          <nav className="hidden md:flex items-center space-x-6">
            <Link href="/" className="text-gray-700 hover:text-primary-600 transition">
              Главная
            </Link>
            <Link href="/news" className="text-gray-700 hover:text-primary-600 transition">
              Новости
            </Link>
            <Link href="/search" className="text-gray-700 hover:text-primary-600 transition">
              Поиск
            </Link>
          </nav>
        </div>
      </div>
    </header>
  )
}
