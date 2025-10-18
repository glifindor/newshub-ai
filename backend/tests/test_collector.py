# 🧪 Pytest тесты для NewsHub AI

import pytest
import pytest_asyncio
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession, async_sessionmaker
from sqlalchemy.pool import NullPool

from app.core.database import Base
from app.models import NewsSource, NewsItem, NewsChannel, SourceType, NewsStatus
from app.services.collector import NewsCollector

# Test database URL
TEST_DATABASE_URL = "postgresql+asyncpg://test:test@localhost:5432/test_newshub"


@pytest_asyncio.fixture
async def test_db():
    """Create test database session"""
    engine = create_async_engine(TEST_DATABASE_URL, poolclass=NullPool)
    
    # Create tables
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    
    # Create session
    async_session = async_sessionmaker(engine, class_=AsyncSession, expire_on_commit=False)
    
    async with async_session() as session:
        yield session
    
    # Drop tables
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.drop_all)
    
    await engine.dispose()


@pytest.mark.asyncio
async def test_content_hash(test_db):
    """Test content hashing for deduplication"""
    collector = NewsCollector(test_db)
    
    title = "Test News Title"
    content = "Test news content"
    
    hash1 = collector.calculate_content_hash(title, content)
    hash2 = collector.calculate_content_hash(title, content)
    
    assert hash1 == hash2
    assert len(hash1) == 32  # MD5 hash length


@pytest.mark.asyncio
async def test_keyword_filtering(test_db):
    """Test news filtering by keywords"""
    collector = NewsCollector(test_db)
    
    # Crypto keywords
    crypto_text = "Bitcoin price reached $100,000 today"
    assert collector.filter_by_keywords(crypto_text, NewsChannel.CRYPTO) == True
    
    # Politics keywords
    politics_text = "Kremlin announced new policy"
    assert collector.filter_by_keywords(politics_text, NewsChannel.POLITICS) == True
    
    # Wrong category
    assert collector.filter_by_keywords(crypto_text, NewsChannel.POLITICS) == False


@pytest.mark.asyncio
async def test_create_news_item(test_db):
    """Test creating a news item"""
    news = NewsItem(
        source_id=1,
        title="Test News",
        content="Test content",
        url="https://example.com/test",
        category=NewsChannel.CRYPTO,
        content_hash="test123",
        status=NewsStatus.PENDING,
    )
    
    test_db.add(news)
    await test_db.commit()
    
    assert news.id is not None
    assert news.status == NewsStatus.PENDING


@pytest.mark.asyncio
async def test_duplicate_detection(test_db):
    """Test duplicate news detection"""
    collector = NewsCollector(test_db)
    
    # Create first news
    news1 = NewsItem(
        source_id=1,
        title="Duplicate Test",
        content="Duplicate content",
        url="https://example.com/test1",
        category=NewsChannel.CRYPTO,
        content_hash="duplicate123",
        status=NewsStatus.PENDING,
    )
    test_db.add(news1)
    await test_db.commit()
    
    # Check for duplicate
    is_duplicate = await collector.check_duplicate("duplicate123")
    assert is_duplicate == True
    
    # Check non-duplicate
    is_not_duplicate = await collector.check_duplicate("unique456")
    assert is_not_duplicate == False


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
