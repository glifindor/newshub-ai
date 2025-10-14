-- Создание отдельных баз данных для каждого сервиса

-- Auth Service Database
CREATE DATABASE auth_db;
\c auth_db;
-- Схема создастся через миграции auth-service

-- News Service Database
CREATE DATABASE news_db;
\c news_db;
-- Схема создастся через миграции news-service

-- SEO Service Database
CREATE DATABASE seo_db;
\c seo_db;
-- Схема создастся через миграции seo-service

-- Media Service Database
CREATE DATABASE media_db;
\c media_db;
-- Схема создастся через миграции media-service

-- Возвращаемся к базе по умолчанию
\c postgres;
