-- =====================================================
-- Migration: Add SEO Optimization Fields to News Table
-- Version: 004
-- Description: Adds dedicated SEO fields (title, description, keywords, image)
-- =====================================================

-- Add SEO fields to news table
ALTER TABLE news 
    ADD COLUMN IF NOT EXISTS seo_title VARCHAR(70),
    ADD COLUMN IF NOT EXISTS seo_description VARCHAR(160),
    ADD COLUMN IF NOT EXISTS seo_keywords VARCHAR(255),
    ADD COLUMN IF NOT EXISTS seo_image VARCHAR(500);

-- Add indexes for better performance on SEO fields
CREATE INDEX IF NOT EXISTS idx_news_seo_title ON news(seo_title) WHERE seo_title IS NOT NULL AND seo_title != '';

-- Add comments for documentation
COMMENT ON COLUMN news.seo_title IS 'Optimized title for search engines (max 60 chars recommended, 70 max)';
COMMENT ON COLUMN news.seo_description IS 'Meta description for search engines (max 160 chars)';
COMMENT ON COLUMN news.seo_keywords IS 'Comma-separated keywords for SEO';
COMMENT ON COLUMN news.seo_image IS 'Image URL for Open Graph and Twitter Cards (social media sharing)';

-- Optional: Initialize SEO fields from existing meta fields for backward compatibility
UPDATE news 
SET 
    seo_title = CASE 
        WHEN meta_title IS NOT NULL AND meta_title != '' THEN LEFT(meta_title, 70)
        ELSE LEFT(title, 60)
    END,
    seo_description = CASE 
        WHEN meta_description IS NOT NULL AND meta_description != '' THEN LEFT(meta_description, 160)
        WHEN summary IS NOT NULL AND summary != '' THEN LEFT(summary, 160)
        ELSE LEFT(content, 160)
    END,
    seo_keywords = CASE 
        WHEN meta_keywords IS NOT NULL AND meta_keywords != '' THEN LEFT(meta_keywords, 255)
        ELSE NULL
    END,
    seo_image = featured_image
WHERE seo_title IS NULL;

-- Success message
DO $$
BEGIN
    RAISE NOTICE 'Migration 004: SEO fields successfully added to news table';
END $$;
