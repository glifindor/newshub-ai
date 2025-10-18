# 🚀 NewsHub AI - Frontend

Современный frontend для агрегатора новостей с AI-анализом на Next.js 14, React 18, TypeScript и Tailwind CSS.

## 📋 Возможности

### Публичная часть:
- 🏠 Главная страница с поиском и фильтрами
- 📰 Карточки новостей с AI-саммари
- 🔐 Архив по каналам (crypto/politics)
- 📱 Responsive дизайн
- 🌓 Темная/светлая тема
- 🔄 Real-time обновления через WebSocket

### Админ-панель:
- 🔐 JWT авторизация через NextAuth
- 📊 Dashboard с таблицей новостей (TanStack Table)
- ✅ Модерация: одобрение/отклонение новостей
- 🗂️ CRUD для источников новостей
- 🔍 Фильтры по каналу/статусу/дате
- 📈 Real-time уведомления

## 🛠️ Технологии

- **Framework:** Next.js 14 (Pages Router)
- **UI:** React 18, TypeScript
- **Styling:** Tailwind CSS
- **State Management:** Zustand, TanStack Query
- **Forms:** React Hook Form + Zod
- **Tables:** TanStack Table
- **Auth:** NextAuth.js
- **API Client:** Axios
- **WebSocket:** Socket.IO Client
- **Animations:** Framer Motion
- **Icons:** React Icons
- **Testing:** Cypress (E2E), Storybook (Components)

## 📦 Установка

### Требования

- Node.js >= 18.0.0
- npm >= 9.0.0
- Backend API запущен на `http://localhost:8000`

### Быстрый старт

```bash
# 1. Создать Next.js проект (если ещё не создан)
npx create-next-app@latest newshub-ai-frontend --typescript --tailwind --eslint
cd newshub-ai-frontend

# 2. Установить зависимости
npm install

# Основные зависимости
npm install next-auth axios socket.io-client @tanstack/react-table @tanstack/react-query
npm install zustand date-fns react-hook-form zod @hookform/resolvers
npm install react-hot-toast framer-motion react-icons clsx tailwind-merge

# Dev зависимости
npm install -D @storybook/react @storybook/addon-essentials @storybook/nextjs
npm install -D cypress start-server-and-test
npm install -D @tailwindcss/typography @tailwindcss/forms

# 3. Скопировать файлы из этого репозитория
# Структура:
# src/
#   ├── pages/
#   ├── components/
#   ├── hooks/
#   ├── lib/
#   ├── types/
#   ├── utils/
#   └── styles/

# 4. Создать .env.local
cp .env.local.example .env.local

# Редактировать .env.local:
NEXT_PUBLIC_API_URL=http://localhost:8000/api/v1
NEXT_PUBLIC_WS_URL=ws://localhost:8000
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-super-secret-key-here

# 5. Запустить dev server
npm run dev
```

Откройте http://localhost:3000

## 📁 Структура проекта

```
frontend/
├── src/
│   ├── pages/                    # Next.js страницы
│   │   ├── index.tsx             # Главная страница
│   │   ├── _app.tsx              # App wrapper
│   │   ├── _document.tsx         # HTML document
│   │   ├── admin/
│   │   │   ├── login.tsx         # Страница логина
│   │   │   ├── dashboard.tsx    # Админ dashboard
│   │   │   └── sources.tsx      # Управление источниками
│   │   ├── public/
│   │   │   └── [channel].tsx    # Архив по каналам
│   │   └── api/
│   │       └── auth/
│   │           └── [...nextauth].ts  # NextAuth API
│   │
│   ├── components/               # React компоненты
│   │   ├── NewsCard.tsx          # Карточка новости
│   │   ├── AdminLayout.tsx       # Layout админки
│   │   ├── PublicLayout.tsx      # Layout публичной части
│   │   ├── NewsTable.tsx         # Таблица новостей (TanStack)
│   │   ├── SourceForm.tsx        # Форма источника
│   │   ├── SearchBar.tsx         # Поиск
│   │   ├── Filters.tsx           # Фильтры
│   │   └── ...                   # Другие компоненты
│   │
│   ├── hooks/                    # Custom hooks
│   │   ├── useNews.ts            # Работа с новостями
│   │   ├── useSources.ts         # Работа с источниками
│   │   ├── useTheme.ts           # Тема
│   │   └── usePipeline.ts        # Pipeline операции
│   │
│   ├── lib/                      # Библиотеки
│   │   ├── api.ts                # API client
│   │   └── websocket.ts          # WebSocket client
│   │
│   ├── types/                    # TypeScript types
│   │   └── index.ts              # Все типы
│   │
│   ├── utils/                    # Утилиты
│   │   └── cn.ts                 # classNames helper
│   │
│   └── styles/                   # Стили
│       └── globals.css           # Global CSS
│
├── public/                       # Статика
│   ├── logo.svg
│   └── ...
│
├── .storybook/                   # Storybook config
├── cypress/                      # E2E тесты
├── package.json
├── tsconfig.json
├── tailwind.config.ts
├── next.config.js
└── README.md
```

## 🎨 Дизайн

### Цветовая схема

```typescript
// Primary (Blue)
primary-500: #0ea5e9
primary-600: #0284c7

// Crypto (Orange)
crypto-500: #f59e0b
crypto-600: #d97706

// Politics (Red)
politics-500: #ef4444
politics-600: #dc2626
```

### Темы

- **Светлая тема:** Белый фон, серые акценты
- **Тёмная тема:** Тёмно-серый фон (#0f172a), светлые акценты

## 🔌 API Integration

### Конфигурация

```typescript
// src/lib/api.ts
const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000/api/v1';

export const apiClient = axios.create({
  baseURL: API_URL,
  timeout: 30000,
});

// Auth interceptor
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

### Основные endpoints

```typescript
// Новости
newsApi.getAll(filters) → GET /news/
newsApi.getById(id) → GET /news/{id}
newsApi.approve(id) → POST /news/{id}/approve
newsApi.reject(id) → POST /news/{id}/reject

// Источники
sourcesApi.getAll() → GET /sources/
sourcesApi.create(data) → POST /sources/
sourcesApi.update(id, data) → PUT /sources/{id}
sourcesApi.delete(id) → DELETE /sources/{id}

// Pipeline
pipelineApi.collect(channel) → POST /pipeline/collect
pipelineApi.analyze(limit) → POST /pipeline/analyze
pipelineApi.post(limit) → POST /pipeline/post
```

## 🔐 Авторизация

### NextAuth.js

```typescript
// src/pages/api/auth/[...nextauth].ts
import NextAuth from 'next-auth';
import CredentialsProvider from 'next-auth/providers/credentials';
import { authApi } from '@/lib/api';

export default NextAuth({
  providers: [
    CredentialsProvider({
      name: 'Credentials',
      credentials: {
        username: { label: "Username", type: "text" },
        password: { label: "Password", type: "password" }
      },
      async authorize(credentials) {
        const tokens = await authApi.login({
          username: credentials!.username,
          password: credentials!.password,
        });
        
        return {
          id: '1',
          email: credentials!.username,
          accessToken: tokens.access_token,
        };
      },
    }),
  ],
  pages: {
    signIn: '/admin/login',
  },
  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        token.accessToken = user.accessToken;
      }
      return token;
    },
    async session({ session, token }) {
      session.accessToken = token.accessToken;
      return session;
    },
  },
});
```

### Использование

```typescript
// В компоненте
import { useSession, signIn, signOut } from 'next-auth/react';

const { data: session, status } = useSession();

if (status === 'loading') return <div>Loading...</div>;
if (!session) return <div>Not logged in</div>;

return <div>Welcome, {session.user.email}!</div>;
```

## 🔄 Real-time Updates

### WebSocket Client

```typescript
// src/lib/websocket.ts
import { io } from 'socket.io-client';

const socket = io('ws://localhost:8000', {
  transports: ['websocket', 'polling'],
});

socket.on('news_published', (data) => {
  console.log('New news published:', data);
  // Refetch queries or update state
});
```

### Использование в компонентах

```typescript
import { useEffect } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { wsClient } from '@/lib/websocket';
import { newsKeys } from '@/hooks/useNews';

function Dashboard() {
  const queryClient = useQueryClient();

  useEffect(() => {
    wsClient.connect();

    wsClient.on('news_published', () => {
      queryClient.invalidateQueries({ queryKey: newsKeys.lists() });
    });

    return () => wsClient.disconnect();
  }, []);

  // ...
}
```

## 🧪 Тестирование

### Unit тесты (Jest)

```bash
npm test
```

### Storybook

```bash
# Запустить Storybook
npm run storybook

# Откроется на http://localhost:6006
```

### E2E тесты (Cypress)

```bash
# Открыть Cypress UI
npm run cypress

# Запустить headless
npm run cypress:headless

# Запустить с dev server
npm run test:e2e
```

## 📚 Примеры использования

### Fetch новостей с фильтрами

```typescript
import { useNews } from '@/hooks/useNews';

function NewsPage() {
  const { data, isLoading } = useNews({
    category: 'crypto',
    status: 'published',
    page: 1,
    per_page: 20,
  });

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      {data?.items.map((news) => (
        <NewsCard key={news.id} news={news} />
      ))}
    </div>
  );
}
```

### Модерация новости

```typescript
import { useApproveNews, useRejectNews } from '@/hooks/useNews';

function ModerationPanel({ newsId }: { newsId: string }) {
  const approveMutation = useApproveNews();
  const rejectMutation = useRejectNews();

  return (
    <div>
      <button onClick={() => approveMutation.mutate(newsId)}>
        ✅ Одобрить
      </button>
      <button onClick={() => rejectMutation.mutate(newsId)}>
        ❌ Отклонить
      </button>
    </div>
  );
}
```

### Создание источника

```typescript
import { useCreateSource } from '@/hooks/useSources';
import { useForm } from 'react-hook-form';

function SourceForm() {
  const { register, handleSubmit } = useForm();
  const createMutation = useCreateSource();

  const onSubmit = (data) => {
    createMutation.mutate(data);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <input {...register('name')} placeholder="Название" />
      <input {...register('url')} placeholder="URL" />
      <button type="submit">Создать</button>
    </form>
  );
}
```

## 🚀 Deployment

### Build для production

```bash
npm run build
npm start
```

### Vercel (рекомендуется для Next.js)

```bash
# Установить Vercel CLI
npm i -g vercel

# Deploy
vercel

# Production deploy
vercel --prod
```

### Docker

```dockerfile
FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

EXPOSE 3000

CMD ["npm", "start"]
```

```bash
docker build -t newshub-frontend .
docker run -p 3000:3000 newshub-frontend
```

## 🔧 Конфигурация

### Environment Variables

```bash
# .env.local

# Backend API
NEXT_PUBLIC_API_URL=http://localhost:8000/api/v1
NEXT_PUBLIC_WS_URL=ws://localhost:8000

# NextAuth
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-super-secret-key-change-this-in-production

# Опционально
NEXT_PUBLIC_GA_ID=G-XXXXXXXXXX  # Google Analytics
```

## 📝 TODO

- [ ] Добавить unit тесты для всех компонентов
- [ ] Настроить CI/CD pipeline
- [ ] Добавить i18n (интернационализация)
- [ ] Оптимизировать bundle size
- [ ] Добавить PWA support
- [ ] Интегрировать Sentry для мониторинга ошибок

## 🤝 Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

MIT

## 👥 Authors

- **NewsHub AI Team** - [GitHub](https://github.com/glifindor/newsportal)

## 🙏 Acknowledgments

- Next.js Team
- Tailwind CSS
- TanStack Team
- NextAuth.js
