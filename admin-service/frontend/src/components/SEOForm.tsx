import { useEffect, useState } from 'react'
import { seoAPI } from '../api/api'

interface SEOFormProps {
  newsId: string
}

export default function SEOForm({ newsId }: SEOFormProps) {
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  const [metaTitle, setMetaTitle] = useState('')
  const [metaDescription, setMetaDescription] = useState('')
  const [metaKeywords, setMetaKeywords] = useState('')
  const [ogTitle, setOgTitle] = useState('')
  const [ogDescription, setOgDescription] = useState('')
  const [ogImage, setOgImage] = useState('')

  useEffect(() => {
    loadSEO()
  }, [newsId])

  const loadSEO = async () => {
    try {
      const data = await seoAPI.getNewsSEO(newsId)
      if (data) {
        setMetaTitle(data.title || '')
        setMetaDescription(data.description || '')
        setMetaKeywords(data.keywords || '')
        setOgTitle(data.og_title || '')
        setOgDescription(data.og_description || '')
        setOgImage(data.og_image || '')
      }
    } catch (error) {
      console.error('Failed to load SEO:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      await seoAPI.updateNewsSEO(newsId, {
        title: metaTitle || undefined,
        description: metaDescription || undefined,
        keywords: metaKeywords || undefined,
        og_title: ogTitle || undefined,
        og_description: ogDescription || undefined,
        og_image: ogImage || undefined,
      })
      alert('SEO-данные обновлены')
    } catch (error) {
      console.error('Failed to save SEO:', error)
      alert('Ошибка сохранения SEO-данных')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return <div className="text-center py-4">Загрузка SEO...</div>
  }

  return (
    <div className="mt-4 space-y-4 p-4 bg-gray-50 rounded-md">
      <h3 className="font-bold text-lg">Метаданные для SEO</h3>

      {/* Meta Title */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Meta Title {metaTitle.length > 0 && `(${metaTitle.length}/70)`}
        </label>
        <input
          type="text"
          value={metaTitle}
          onChange={(e) => setMetaTitle(e.target.value)}
          maxLength={70}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
          placeholder="Заголовок для поисковиков"
        />
      </div>

      {/* Meta Description */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Meta Description {metaDescription.length > 0 && `(${metaDescription.length}/160)`}
        </label>
        <textarea
          value={metaDescription}
          onChange={(e) => setMetaDescription(e.target.value)}
          maxLength={160}
          rows={3}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
          placeholder="Описание для поисковиков"
        />
      </div>

      {/* Meta Keywords */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Keywords (через запятую)
        </label>
        <input
          type="text"
          value={metaKeywords}
          onChange={(e) => setMetaKeywords(e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
          placeholder="новости, политика, экономика"
        />
      </div>

      <h3 className="font-bold text-lg mt-6">Open Graph (социальные сети)</h3>

      {/* OG Title */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          OG Title
        </label>
        <input
          type="text"
          value={ogTitle}
          onChange={(e) => setOgTitle(e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
          placeholder="Заголовок для социальных сетей"
        />
      </div>

      {/* OG Description */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          OG Description
        </label>
        <textarea
          value={ogDescription}
          onChange={(e) => setOgDescription(e.target.value)}
          rows={2}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
          placeholder="Описание для социальных сетей"
        />
      </div>

      {/* OG Image */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          OG Image (URL)
        </label>
        <input
          type="text"
          value={ogImage}
          onChange={(e) => setOgImage(e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
          placeholder="https://example.com/og-image.jpg"
        />
        {ogImage && (
          <img src={ogImage} alt="OG Preview" className="mt-2 max-w-xs rounded border" />
        )}
      </div>

      {/* Save Button */}
      <button
        onClick={handleSave}
        disabled={saving}
        className="px-6 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-md disabled:opacity-50"
      >
        {saving ? 'Сохранение...' : 'Сохранить SEO'}
      </button>
    </div>
  )
}
