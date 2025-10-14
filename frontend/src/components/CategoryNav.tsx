'use client'

import Link from 'next/link'
import { useState, useEffect } from 'react'
import { getCategories } from '@/lib/api'

export default function CategoryNav() {
  const [categories, setCategories] = useState<string[]>([])

  useEffect(() => {
    getCategories().then(setCategories)
  }, [])

  return (
    <nav className="bg-white border-b sticky top-16 z-40">
      <div className="container mx-auto px-4">
        <ul className="flex items-center gap-1 overflow-x-auto py-3">
          <li>
            <Link
              href="/news"
              className="px-4 py-2 rounded-lg hover:bg-gray-100 text-gray-700 font-medium whitespace-nowrap transition"
            >
              Все
            </Link>
          </li>
          {categories.map((category) => (
            <li key={category}>
              <Link
                href={`/category/${encodeURIComponent(category)}`}
                className="px-4 py-2 rounded-lg hover:bg-primary-50 hover:text-primary-700 text-gray-700 whitespace-nowrap transition"
              >
                {category}
              </Link>
            </li>
          ))}
        </ul>
      </div>
    </nav>
  )
}
