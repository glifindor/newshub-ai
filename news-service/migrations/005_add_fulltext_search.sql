-- =====================================================
-- Migration: Add Full Text Search Support
-- Version: 005
-- Description: Adds tsvector column, triggers, and GIN index for PostgreSQL Full Text Search
-- =====================================================

-- Step 1: Add tsvector column for full-text search
ALTER TABLE news ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- Step 2: Create function to automatically update search vector
-- This function assigns different weights to different fields:
-- 'A' (highest) - title, seo_title
-- 'B' (high) - summary (excerpt), seo_description
-- 'C' (medium) - content
-- 'D' (low) - seo_keywords
CREATE OR REPLACE FUNCTION news_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.seo_title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.summary, '')), 'B') ||
        setweight(to_tsvector('english', coalesce(NEW.seo_description, '')), 'B') ||
        setweight(to_tsvector('english', coalesce(NEW.content, '')), 'C') ||
        setweight(to_tsvector('english', coalesce(NEW.seo_keywords, '')), 'D');
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

-- Step 3: Create trigger for automatic updates on INSERT and UPDATE
DROP TRIGGER IF EXISTS news_search_vector_trigger ON news;
CREATE TRIGGER news_search_vector_trigger
    BEFORE INSERT OR UPDATE OF title, content, summary, seo_title, seo_description, seo_keywords
    ON news
    FOR EACH ROW
    EXECUTE FUNCTION news_search_vector_update();

-- Step 4: Populate search_vector for existing rows
UPDATE news SET search_vector = 
    setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(seo_title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(summary, '')), 'B') ||
    setweight(to_tsvector('english', coalesce(seo_description, '')), 'B') ||
    setweight(to_tsvector('english', coalesce(content, '')), 'C') ||
    setweight(to_tsvector('english', coalesce(seo_keywords, '')), 'D')
WHERE search_vector IS NULL OR search_vector = '';

-- Step 5: Create GIN index for fast full-text searching
-- GIN (Generalized Inverted Index) is optimized for full-text search
CREATE INDEX IF NOT EXISTS idx_news_search_vector ON news USING GIN(search_vector);

-- Step 6: Add index on published status for better query performance
CREATE INDEX IF NOT EXISTS idx_news_status_published ON news(status, published_at DESC) 
WHERE status = 'published';

-- Step 7: Create helper function for ranked search with results
-- This function returns news with relevance ranking
CREATE OR REPLACE FUNCTION news_search_ranked(search_query text, max_results integer DEFAULT 100)
RETURNS TABLE (
    id uuid,
    title varchar,
    summary text,
    slug varchar,
    published_at timestamp,
    rank real
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        n.id,
        n.title,
        n.summary,
        n.slug,
        n.published_at,
        ts_rank(n.search_vector, plainto_tsquery('english', search_query)) as rank
    FROM news n
    WHERE n.search_vector @@ plainto_tsquery('english', search_query)
        AND n.status = 'published'
    ORDER BY rank DESC, n.published_at DESC
    LIMIT max_results;
END;
$$ LANGUAGE plpgsql STABLE;

-- Step 8: Create function for search with headline (snippet) generation
CREATE OR REPLACE FUNCTION news_search_with_snippet(search_query text, max_results integer DEFAULT 100)
RETURNS TABLE (
    id uuid,
    title varchar,
    slug varchar,
    snippet text,
    rank real
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        n.id,
        n.title,
        n.slug,
        ts_headline('english', n.content, plainto_tsquery('english', search_query),
            'StartSel=<mark>, StopSel=</mark>, MaxWords=50, MinWords=20') as snippet,
        ts_rank(n.search_vector, plainto_tsquery('english', search_query)) as rank
    FROM news n
    WHERE n.search_vector @@ plainto_tsquery('english', search_query)
        AND n.status = 'published'
    ORDER BY rank DESC, n.published_at DESC
    LIMIT max_results;
END;
$$ LANGUAGE plpgsql STABLE;

-- Add comments for documentation
COMMENT ON COLUMN news.search_vector IS 'Automatically updated tsvector for full-text search (weighted: A=title, B=summary, C=content, D=keywords)';
COMMENT ON FUNCTION news_search_vector_update() IS 'Trigger function to automatically update search_vector when news content changes';
COMMENT ON FUNCTION news_search_ranked(text, integer) IS 'Search published news with relevance ranking';
COMMENT ON FUNCTION news_search_with_snippet(text, integer) IS 'Search published news with highlighted snippets';

-- Success message
DO $$
BEGIN
    RAISE NOTICE 'Migration 005: Full Text Search successfully configured';
    RAISE NOTICE '  - search_vector column added with automatic trigger updates';
    RAISE NOTICE '  - GIN index created for fast searching';
    RAISE NOTICE '  - Helper functions: news_search_ranked() and news_search_with_snippet()';
END $$;
