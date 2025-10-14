import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { newsAPI, CreateNewsRequest, UpdateNewsRequest } from '../api/api'
import TipTapEditor from '../components/TipTapEditor'
import SEOForm from '../components/SEOForm'

export default function NewsEdit() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const isEditMode = !!id

  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)

  const [title, setTitle] = useState('')
  const [slug, setSlug] = useState('')
  const [summary, setSummary] = useState('')
  const [content, setContent] = useState('')
  const [category, setCategory] = useState('')
  const [tags, setTags] = useState<string[]>([])
  const [tagInput, setTagInput] = useState('')
  const [coverImage, setCoverImage] = useState('')
  const [showSEO, setShowSEO] = useState(false)

  useEffect(() => {
    if (isEditMode && id) {
      loadNews(id)
    }
  }, [id])

  const loadNews = async (newsId: string) => {
    setLoading(true)
    try {
      const data = await newsAPI.getNewsById(newsId)
      setTitle(data.title)
      setSlug(data.slug)
      setSummary(data.summary || '')
      setContent(data.content)
      setCategory(data.category || '')
      setTags(data.tags || [])
      setCoverImage(data.image_url || '')
    } catch (error) {
      console.error('Failed to load news:', error)
      alert('Ошибка загрузки новости')
    } finally {
      setLoading(false)
    }
  }

  const generateSlug = (text: string) => {
    return text
      .toLowerCase()
      .replace(/[^\w\s-]/g, '')
      .replace(/\s+/g, '-')
      .replace(/-+/g, '-')
      .trim()
  }

  const handleTitleChange = (value: string) => {
    setTitle(value)
    if (!isEditMode && !slug) {
      setSlug(generateSlug(value))
    }
  }

  const handleAddTag = () => {
    if (tagInput && !tags.includes(tagInput)) {
      setTags([...tags, tagInput])
      setTagInput('')
    }
  }

  const handleRemoveTag = (tag: string) => {
    setTags(tags.filter((t) => t !== tag))
  }

  const handleSave = async (publishStatus: 'draft' | 'published') => {
    if (!title || !content) {
      alert('Заполните заголовок и содержание')
      return
    }

    setSaving(true)
    try {
      const data: CreateNewsRequest | UpdateNewsRequest = {
        title,
        slug,
        summary,
        content,
        tags: tags.length > 0 ? tags : undefined,
        image_url: coverImage || undefined,
        status: publishStatus,
      }

      if (isEditMode && id) {
        await newsAPI.updateNews(id, data as UpdateNewsRequest)
        alert('Новость обновлена')
      } else {
        const created = await newsAPI.createNews(data as CreateNewsRequest)
        alert('Новость создана')
        navigate(`/news/${created.id}`)
      }
    } catch (error) {
      console.error('Failed to save news:', error)
      alert('Ошибка сохранения новости')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return <div className="text-center py-12">Загрузка...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">
          {isEditMode ? 'Редактирование новости' : 'Создание новости'}
        </h1>
        <button
          onClick={() => navigate('/news')}
          className="text-gray-600 hover:text-gray-900"
        >
          ← Назад к списку
        </button>
      </div>

      <div className="bg-white rounded-lg shadow p-6 space-y-6">
        {/* Title */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Заголовок *
          </label>
          <input
            type="text"
            value={title}
            onChange={(e) => handleTitleChange(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
            placeholder="Введите заголовок новости"
          />
        </div>

        {/* Slug */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            URL (slug) *
          </label>
          <input
            type="text"
            value={slug}
            onChange={(e) => setSlug(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
            placeholder="url-novosti"
          />
        </div>

        {/* Summary */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Краткое описание
          </label>
          <textarea
            value={summary}
            onChange={(e) => setSummary(e.target.value)}
            rows={3}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
            placeholder="Краткое описание новости"
          />
        </div>

        {/* Content Editor */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Содержание *
          </label>
          <TipTapEditor content={content} onChange={setContent} />
        </div>

        {/* Category & Tags */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Категория
            </label>
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
            >
              <option value="">Выберите категорию</option>
              <option value="Политика">Политика</option>
              <option value="Экономика">Экономика</option>
              <option value="Технологии">Технологии</option>
              <option value="Спорт">Спорт</option>
              <option value="Культура">Культура</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Теги
            </label>
            <div className="flex gap-2">
              <input
                type="text"
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), handleAddTag())}
                className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
                placeholder="Добавить тег"
              />
              <button
                type="button"
                onClick={handleAddTag}
                className="px-4 py-2 bg-gray-200 hover:bg-gray-300 rounded-md"
              >
                +
              </button>
            </div>
            <div className="flex flex-wrap gap-2 mt-2">
              {tags.map((tag) => (
                <span
                  key={tag}
                  className="px-3 py-1 bg-primary-100 text-primary-700 rounded-full text-sm flex items-center gap-2"
                >
                  {tag}
                  <button
                    onClick={() => handleRemoveTag(tag)}
                    className="text-primary-900 hover:text-primary-600"
                  >
                    ×
                  </button>
                </span>
              ))}
            </div>
          </div>
        </div>

        {/* Cover Image */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Обложка (URL)
          </label>
          <input
            type="text"
            value={coverImage}
            onChange={(e) => setCoverImage(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
            placeholder="https://example.com/image.jpg"
          />
          {coverImage && (
            <img src={coverImage} alt="Preview" className="mt-2 max-w-xs rounded" />
          )}
        </div>

        {/* SEO Section */}
        <div className="border-t pt-6">
          <button
            type="button"
            onClick={() => setShowSEO(!showSEO)}
            className="flex items-center gap-2 text-gray-700 hover:text-gray-900 font-medium"
          >
            {showSEO ? '▼' : '▶'} SEO-настройки
          </button>
          {showSEO && isEditMode && id && <SEOForm newsId={id} />}
        </div>

        {/* Actions */}
        <div className="flex gap-4 pt-6 border-t">
          <button
            onClick={() => handleSave('draft')}
            disabled={saving}
            className="px-6 py-2 border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50"
          >
            Сохранить черновик
          </button>
          <button
            onClick={() => handleSave('published')}
            disabled={saving}
            className="px-6 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-md disabled:opacity-50"
          >
            {isEditMode ? 'Обновить и опубликовать' : 'Опубликовать'}
          </button>
        </div>
      </div>
    </div>
  )
}
