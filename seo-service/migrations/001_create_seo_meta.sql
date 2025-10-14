-- =====================================================
-- Migration: Create SEO Meta Table
-- Version: 001
-- Description: Таблица для хранения SEO метаданных новостей
-- =====================================================

-- Создание таблицы seo_meta
CREATE TABLE IF NOT EXISTS seo_meta (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    news_id UUID NOT NULL UNIQUE,
    
    -- Основные SEO поля
    title VARCHAR(70) NOT NULL,
    description VARCHAR(160) NOT NULL,
    keywords VARCHAR(255),
    slug VARCHAR(500) NOT NULL UNIQUE,
    canonical_url VARCHAR(500),
    
    -- Open Graph (Facebook, LinkedIn, etc.)
    og_title VARCHAR(100),
    og_description VARCHAR(200),
    og_image VARCHAR(500),
    og_type VARCHAR(50) DEFAULT 'article',
    og_locale VARCHAR(10) DEFAULT 'en_US',
    
    -- Twitter Cards
    twitter_card VARCHAR(50) DEFAULT 'summary_large_image',
    twitter_title VARCHAR(100),
    twitter_description VARCHAR(200),
    twitter_image VARCHAR(500),
    twitter_creator VARCHAR(50),
    
    -- Schema.org (Structured Data)
    schema_type VARCHAR(50) DEFAULT 'NewsArticle',
    schema_data JSONB,
    
    -- Indexing & Sitemap
    robots_index BOOLEAN DEFAULT true,
    robots_follow BOOLEAN DEFAULT true,
    sitemap_priority DECIMAL(2,1) DEFAULT 0.5 CHECK (sitemap_priority >= 0.0 AND sitemap_priority <= 1.0),
    sitemap_changefreq VARCHAR(20) DEFAULT 'weekly' CHECK (sitemap_changefreq IN ('always', 'hourly', 'daily', 'weekly', 'monthly', 'yearly', 'never')),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign Key (если таблица news существует)
    CONSTRAINT fk_seo_meta_news FOREIGN KEY (news_id) 
        REFERENCES news(id) 
        ON DELETE CASCADE
);

-- Индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_seo_meta_slug ON seo_meta(slug);
CREATE INDEX IF NOT EXISTS idx_seo_meta_news_id ON seo_meta(news_id);
CREATE INDEX IF NOT EXISTS idx_seo_meta_updated_at ON seo_meta(updated_at DESC);

-- Индекс для sitemap (только индексируемые записи)
CREATE INDEX IF NOT EXISTS idx_seo_meta_sitemap 
    ON seo_meta(sitemap_priority DESC, updated_at DESC) 
    WHERE robots_index = true;

-- Индекс для JSONB поля (для быстрого поиска по schema_data)
CREATE INDEX IF NOT EXISTS idx_seo_meta_schema_data 
    ON seo_meta USING GIN(schema_data);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_seo_meta_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS seo_meta_updated_at_trigger ON seo_meta;
CREATE TRIGGER seo_meta_updated_at_trigger
    BEFORE UPDATE ON seo_meta
    FOR EACH ROW
    EXECUTE FUNCTION update_seo_meta_updated_at();

-- Комментарии к таблице и колонкам
COMMENT ON TABLE seo_meta IS 'SEO метаданные для новостей (Open Graph, Twitter Cards, Schema.org)';

COMMENT ON COLUMN seo_meta.id IS 'Уникальный идентификатор SEO записи';
COMMENT ON COLUMN seo_meta.news_id IS 'ID новости (foreign key на news.id)';
COMMENT ON COLUMN seo_meta.title IS 'SEO title (рекомендуется 50-60 символов, макс 70)';
COMMENT ON COLUMN seo_meta.description IS 'Meta description (рекомендуется 150-160 символов)';
COMMENT ON COLUMN seo_meta.keywords IS 'Ключевые слова через запятую (устаревшее, но может использоваться)';
COMMENT ON COLUMN seo_meta.slug IS 'URL slug новости (уникальный)';
COMMENT ON COLUMN seo_meta.canonical_url IS 'Canonical URL для избежания дублирования контента';

COMMENT ON COLUMN seo_meta.og_title IS 'Open Graph title (может отличаться от SEO title)';
COMMENT ON COLUMN seo_meta.og_description IS 'Open Graph description';
COMMENT ON COLUMN seo_meta.og_image IS 'Open Graph image URL (рекомендуется 1200x630px)';
COMMENT ON COLUMN seo_meta.og_type IS 'Open Graph type (article, website, video, etc.)';
COMMENT ON COLUMN seo_meta.og_locale IS 'Локаль контента (en_US, ru_RU, etc.)';

COMMENT ON COLUMN seo_meta.twitter_card IS 'Тип Twitter Card (summary, summary_large_image, player, app)';
COMMENT ON COLUMN seo_meta.twitter_title IS 'Twitter title';
COMMENT ON COLUMN seo_meta.twitter_description IS 'Twitter description';
COMMENT ON COLUMN seo_meta.twitter_image IS 'Twitter image URL (рекомендуется 1200x628px)';
COMMENT ON COLUMN seo_meta.twitter_creator IS 'Twitter @username автора';

COMMENT ON COLUMN seo_meta.schema_type IS 'Тип Schema.org (NewsArticle, BlogPosting, Article, etc.)';
COMMENT ON COLUMN seo_meta.schema_data IS 'Полные данные Schema.org в формате JSON-LD';

COMMENT ON COLUMN seo_meta.robots_index IS 'Разрешить индексацию поисковыми системами (true/false)';
COMMENT ON COLUMN seo_meta.robots_follow IS 'Разрешить переход по ссылкам (true/false)';
COMMENT ON COLUMN seo_meta.sitemap_priority IS 'Приоритет в sitemap.xml (0.0 - 1.0)';
COMMENT ON COLUMN seo_meta.sitemap_changefreq IS 'Частота изменения для sitemap.xml';

COMMENT ON COLUMN seo_meta.created_at IS 'Дата создания записи';
COMMENT ON COLUMN seo_meta.updated_at IS 'Дата последнего обновления (автоматически)';

-- Функция для валидации slug (только латиница, цифры, дефис)
CREATE OR REPLACE FUNCTION validate_slug(slug_value VARCHAR)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN slug_value ~ '^[a-z0-9]+(?:-[a-z0-9]+)*$';
END;
$$ LANGUAGE plpgsql;

-- Constraint для валидации slug
ALTER TABLE seo_meta 
    ADD CONSTRAINT check_slug_format 
    CHECK (validate_slug(slug));

-- Успешное сообщение
DO $$
BEGIN
    RAISE NOTICE 'Migration 001: SEO Meta table created successfully';
    RAISE NOTICE '  - Table: seo_meta';
    RAISE NOTICE '  - Indexes: 5 (slug, news_id, updated_at, sitemap, schema_data)';
    RAISE NOTICE '  - Triggers: 1 (auto-update updated_at)';
    RAISE NOTICE '  - Constraints: 1 (slug format validation)';
END $$;
