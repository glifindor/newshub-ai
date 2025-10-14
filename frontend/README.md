# Frontend для публичного сайта

Next.js 14 приложение для пользователей новостного портала.

## Возможности

- ✅ Главная страница с последними новостями
- ✅ Детальная страница новости с SEO
- ✅ Категории и фильтрация
- ✅ Responsive дизайн
- ✅ Tailwind CSS
- ✅ SSR (Server-Side Rendering)
- ✅ ISR (Incremental Static Regeneration)
- ✅ Счетчик просмотров
- ✅ Теги и категории

## Установка

```bash
npm install
```

## Запуск

### Development
```bash
npm run dev
```

### Production
```bash
npm run build
npm start
```

## Переменные окружения

Создайте `.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Docker

```bash
docker build -t news-frontend .
docker run -p 3000:3000 news-frontend
```

## Структура

```
frontend/
├── src/
│   ├── app/
│   │   ├── layout.tsx         # Root layout
│   │   ├── page.tsx           # Homepage
│   │   ├── news/
│   │   │   └── [slug]/
│   │   │       └── page.tsx   # News detail page
│   │   └── globals.css        # Global styles
│   ├── components/
│   │   ├── Header.tsx         # Header component
│   │   ├── Footer.tsx         # Footer component
│   │   ├── NewsGrid.tsx       # News grid
│   │   └── CategoryNav.tsx    # Category navigation
│   └── lib/
│       └── api.ts             # API client
├── public/                    # Static files
├── next.config.js             # Next.js config
├── tailwind.config.ts         # Tailwind config
├── Dockerfile                 # Docker build
└── package.json
```

## API Endpoints

Приложение использует следующие API:

- `GET /api/v1/news` - Список новостей
- `GET /api/v1/news/slug/:slug` - Новость по slug
- `POST /api/v1/news/:id/view` - Увеличить счетчик
- `GET /api/v1/categories` - Список категорий
