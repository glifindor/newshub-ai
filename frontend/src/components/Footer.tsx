export default function Footer() {
  return (
    <footer className="bg-gray-900 text-white mt-16">
      <div className="container mx-auto px-4 py-12">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div>
            <h3 className="text-xl font-bold mb-4">Новостной портал</h3>
            <p className="text-gray-400">
              Актуальные новости и события из разных областей жизни.
            </p>
          </div>
          
          <div>
            <h4 className="font-bold mb-4">Разделы</h4>
            <ul className="space-y-2 text-gray-400">
              <li><a href="/category/политика" className="hover:text-white transition">Политика</a></li>
              <li><a href="/category/экономика" className="hover:text-white transition">Экономика</a></li>
              <li><a href="/category/технологии" className="hover:text-white transition">Технологии</a></li>
              <li><a href="/category/спорт" className="hover:text-white transition">Спорт</a></li>
            </ul>
          </div>
          
          <div>
            <h4 className="font-bold mb-4">Контакты</h4>
            <ul className="space-y-2 text-gray-400">
              <li>Email: info@newsportal.ru</li>
              <li>Телефон: +7 (xxx) xxx-xx-xx</li>
            </ul>
          </div>
        </div>
        
        <div className="border-t border-gray-800 mt-8 pt-8 text-center text-gray-400">
          <p>&copy; {new Date().getFullYear()} Новостной портал. Все права защищены.</p>
        </div>
      </div>
    </footer>
  )
}
