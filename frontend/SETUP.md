# 🎯 FRONTEND SETUP - Быстрый старт

## ⚡ Установка за 5 минут

### 1. Создать проект

```powershell
# Перейти в корень
cd "C:\Users\Grifindor\Desktop\НОВСТНОЙ САЙТ"

# Создать Next.js проект
npx create-next-app@latest frontend --typescript --tailwind --eslint --src-dir --app=false --import-alias="@/*"

cd frontend
```

При запросе выберите:
- ✅ TypeScript
- ✅ ESLint
- ✅ Tailwind CSS
- ✅ `src/` directory
- ❌ App Router (используем Pages Router)
- ✅ Import alias (@/*)

### 2. Установить зависимости

```powershell
# Основные
npm install next-auth@latest axios socket.io-client
npm install @tanstack/react-table@latest @tanstack/react-query@latest
npm install zustand date-fns react-hook-form zod @hookform/resolvers
npm install react-hot-toast framer-motion react-icons
npm install clsx tailwind-merge

# Dev dependencies
npm install -D @tailwindcss/typography @tailwindcss/forms
npm install -D @types/node
```

### 3. Создать .env.local

```powershell
# Скопировать пример
copy .env.local.example .env.local
```

Редактировать `.env.local`:
```env
NEXT_PUBLIC_API_URL=http://localhost:8000/api/v1
NEXT_PUBLIC_WS_URL=ws://localhost:8000
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=ad7f162fd8215a112366b6b06c02562fc36d34a7e875ebc00b857c62802a57bf
```

### 4. Запустить dev server

```powershell
npm run dev
```

Откройте http://localhost:3000

## 📁 Файлы для создания вручную

Все ключевые файлы уже созданы в этом репозитории:

### ✅ Уже созданы:

1. `package.json` - зависимости
2. `tsconfig.json` - TypeScript конфигурация
3. `tailwind.config.ts` - Tailwind CSS конфигурация
4. `.env.local.example` - пример environment variables
5. `src/types/index.ts` - TypeScript типы
6. `src/lib/api.ts` - API client
7. `src/lib/websocket.ts` - WebSocket client
8. `src/hooks/useNews.ts` - hooks для новостей
9. `src/hooks/useSources.ts` - hooks для источников
10. `src/hooks/useTheme.ts` - hook для темы
11. `src/components/NewsCard.tsx` - карточка новости
12. `src/components/AdminLayout.tsx` - layout админки
13. `src/utils/cn.ts` - classNames utility

### 📝 Нужно создать:

#### 1. `src/pages/_app.tsx`

```typescript
import '@/styles/globals.css';
import type { AppProps } from 'next/app';
import { SessionProvider } from 'next-auth/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Toaster } from 'react-hot-toast';
import { useEffect } from 'react';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

export default function App({ Component, pageProps: { session, ...pageProps } }: AppProps) {
  // Initialize theme
  useEffect(() => {
    const stored = localStorage.getItem('theme');
    if (stored === 'dark' || (!stored && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
      document.documentElement.classList.add('dark');
    }
  }, []);

  return (
    <SessionProvider session={session}>
      <QueryClientProvider client={queryClient}>
        <Component {...pageProps} />
        <Toaster position="top-right" />
      </QueryClientProvider>
    </SessionProvider>
  );
}
```

#### 2. `src/pages/index.tsx` - Главная страница

```typescript
import { useState } from 'react';
import Head from 'next/head';
import Link from 'next/link';
import { usePublicNews } from '@/hooks/useNews';
import NewsCard from '@/components/NewsCard';
import { FiSearch } from 'react-icons/fi';

export default function Home() {
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  
  const { data, isLoading } = usePublicNews({
    search,
    page,
    per_page: 12,
  });

  return (
    <>
      <Head>
        <title>NewsHub AI - Новости с AI-анализом</title>
      </Head>

      <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
        {/* Header */}
        <header className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
          <div className="max-w-7xl mx-auto px-4 py-6">
            <div className="flex items-center justify-between">
              <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
                NewsHub AI 🤖
              </h1>
              <Link
                href="/admin/login"
                className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
              >
                Admin
              </Link>
            </div>

            {/* Search */}
            <div className="mt-6 relative">
              <FiSearch className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
              <input
                type="text"
                placeholder="Поиск новостей..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="w-full pl-10 pr-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 text-gray-900 dark:text-white"
              />
            </div>

            {/* Channel links */}
            <div className="mt-4 flex space-x-4">
              <Link href="/public/crypto" className="text-crypto-600 hover:underline">
                🔐 Crypto
              </Link>
              <Link href="/public/politics" className="text-politics-600 hover:underline">
                🏛️ Politics
              </Link>
            </div>
          </div>
        </header>

        {/* News Grid */}
        <main className="max-w-7xl mx-auto px-4 py-8">
          {isLoading ? (
            <div>Loading...</div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {data?.items.map((news) => (
                <NewsCard key={news.id} news={news} />
              ))}
            </div>
          )}

          {/* Pagination */}
          {data && data.pages > 1 && (
            <div className="mt-8 flex justify-center space-x-2">
              <button
                onClick={() => setPage(p => Math.max(1, p - 1))}
                disabled={page === 1}
                className="px-4 py-2 bg-white dark:bg-gray-800 border rounded-lg disabled:opacity-50"
              >
                Назад
              </button>
              <span className="px-4 py-2">
                {page} / {data.pages}
              </span>
              <button
                onClick={() => setPage(p => Math.min(data.pages, p + 1))}
                disabled={page === data.pages}
                className="px-4 py-2 bg-white dark:bg-gray-800 border rounded-lg disabled:opacity-50"
              >
                Вперёд
              </button>
            </div>
          )}
        </main>
      </div>
    </>
  );
}
```

#### 3. `src/pages/admin/login.tsx`

```typescript
import { useState } from 'react';
import { signIn } from 'next-auth/react';
import { useRouter } from 'next/router';
import { toast } from 'react-hot-toast';

export default function LoginPage() {
  const router = useRouter();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const result = await signIn('credentials', {
        username,
        password,
        redirect: false,
      });

      if (result?.error) {
        toast.error('Неверный логин или пароль');
      } else {
        toast.success('Добро пожаловать!');
        router.push('/admin/dashboard');
      }
    } catch (error) {
      toast.error('Ошибка входа');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
      <div className="max-w-md w-full space-y-8 p-8 bg-white dark:bg-gray-800 rounded-xl shadow-lg">
        <div>
          <h2 className="text-center text-3xl font-bold text-gray-900 dark:text-white">
            NewsHub AI
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
            Админ-панель
          </p>
        </div>

        <form onSubmit={handleSubmit} className="mt-8 space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              Username
            </label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-900 text-gray-900 dark:text-white"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              Password
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-900 text-gray-900 dark:text-white"
              required
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50"
          >
            {loading ? 'Загрузка...' : 'Войти'}
          </button>
        </form>
      </div>
    </div>
  );
}
```

#### 4. `src/styles/globals.css`

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  html {
    @apply antialiased;
  }
  
  body {
    @apply bg-white dark:bg-gray-900 text-gray-900 dark:text-white;
  }
}
```

## ✅ Проверка работы

### 1. Запустить backend

```powershell
cd backend
uvicorn app.main:app --reload
```

Backend должен быть на http://localhost:8000

### 2. Запустить frontend

```powershell
cd frontend
npm run dev
```

Frontend на http://localhost:3000

### 3. Тестирование

- ✅ Открыть http://localhost:3000 - должна загрузиться главная
- ✅ Открыть http://localhost:3000/admin/login - страница логина
- ✅ Войти (admin/password)
- ✅ Должен редиректить на /admin/dashboard

## 🐛 Troubleshooting

### Ошибка: "Cannot find module"

```powershell
rm -rf node_modules
npm install
```

### Ошибка: "API connection failed"

Проверьте что backend запущен:
```powershell
curl http://localhost:8000/api/v1/news/
```

### Ошибка NextAuth

Проверьте `.env.local`:
```env
NEXTAUTH_SECRET=ваш-секретный-ключ
NEXTAUTH_URL=http://localhost:3000
```

## 📚 Следующие шаги

1. Создать остальные страницы (dashboard, sources)
2. Добавить Storybook для компонентов
3. Настроить Cypress для E2E тестов
4. Deploy на Vercel

Документация: `README.md`
