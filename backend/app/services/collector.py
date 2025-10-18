"""
News Collector Service - сбор новостей из RSS и API
"""

import asyncio
import hashlib
from datetime import datetime, timedelta
from typing import List, Optional

import feedparser
import httpx
from sqlalchemy import and_, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.config import settings
from app.core.logger import get_logger
from app.models import NewsChannel, NewsItem, NewsSource, NewsStatus, SourceType

logger = get_logger(__name__)


class NewsCollector:
    """Сборщик новостей из различных источников"""

    def __init__(self, db: AsyncSession):
        self.db = db

        # Keywords для фильтрации
        self.crypto_keywords = [
            "bitcoin",
            "btc",
            "ethereum",
            "eth",
            "crypto",
            "cryptocurrency",
            "blockchain",
            "nft",
            "defi",
            "web3",
            "altcoin",
            "token",
            "binance",
            "coinbase",
            "mining",
            "wallet",
        ]

        self.politics_keywords = [
            "kremlin",
            "putin",
            "russia",
            "ukraine",
            "sanctions",
            "moscow",
            "government",
            "president",
            "minister",
            "duma",
            "foreign policy",
            "diplomacy",
            "geopolitics",
        ]

    def calculate_content_hash(self, title: str, content: str) -> str:
        """Вычислить MD5 хэш для дедупликации"""
        combined = f"{title}|{content}".encode("utf-8")
        return hashlib.md5(combined).hexdigest()

    def filter_by_keywords(self, text: str, channel: NewsChannel) -> bool:
        """Фильтрация по ключевым словам"""
        text_lower = text.lower()

        keywords = (
            self.crypto_keywords
            if channel == NewsChannel.CRYPTO
            else self.politics_keywords
        )

        return any(keyword in text_lower for keyword in keywords)

    async def check_duplicate(self, content_hash: str) -> bool:
        """Проверка на дубликат"""
        result = await self.db.execute(
            select(NewsItem).where(NewsItem.content_hash == content_hash)
        )
        return result.scalar_one_or_none() is not None

    async def parse_rss_feed(self, source: NewsSource) -> List[dict]:
        """Парсинг RSS фида"""
        logger.info(f"Parsing RSS feed: {source.name}", source_id=source.id)

        try:
            # Парсинг RSS (feedparser работает синхронно)
            feed = await asyncio.to_thread(feedparser.parse, source.url)

            if feed.bozo:
                logger.warning(
                    f"RSS feed has errors: {source.name}",
                    error=str(feed.bozo_exception),
                )

            news_items = []

            for entry in feed.entries[: settings.MAX_NEWS_PER_SOURCE]:
                # Извлечение данных
                title = entry.get("title", "")
                content = entry.get("summary", entry.get("description", ""))
                url = entry.get("link", "")
                author = entry.get("author", "")

                # Дата публикации
                published_at = None
                if hasattr(entry, "published_parsed") and entry.published_parsed:
                    published_at = datetime(*entry.published_parsed[:6])

                # Изображение
                image_url = None
                if hasattr(entry, "media_content") and entry.media_content:
                    image_url = entry.media_content[0].get("url")
                elif hasattr(entry, "media_thumbnail") and entry.media_thumbnail:
                    image_url = entry.media_thumbnail[0].get("url")

                # Фильтрация по ключевым словам
                if not self.filter_by_keywords(f"{title} {content}", source.category):
                    continue

                # Проверка на дубликат
                content_hash = self.calculate_content_hash(title, content)
                if await self.check_duplicate(content_hash):
                    logger.debug(f"Duplicate found: {title[:50]}...")
                    continue

                news_items.append(
                    {
                        "title": title,
                        "content": content,
                        "url": url,
                        "author": author,
                        "image_url": image_url,
                        "published_at": published_at,
                        "content_hash": content_hash,
                    }
                )

            logger.info(f"Collected {len(news_items)} news from {source.name}")
            return news_items

        except Exception as e:
            logger.error(f"Error parsing RSS feed {source.name}: {e}", exc_info=True)

            # Обновить счётчик ошибок
            source.error_count += 1
            source.last_error = str(e)
            await self.db.commit()

            return []

    async def fetch_from_newsapi(self, source: NewsSource) -> List[dict]:
        """Получение новостей через NewsAPI.org"""
        logger.info(f"Fetching from NewsAPI: {source.name}")

        try:
            # Определение query в зависимости от категории
            if source.category == NewsChannel.CRYPTO:
                query = "cryptocurrency OR bitcoin OR ethereum OR blockchain"
            else:
                query = "russia OR kremlin OR putin OR ukraine"

            # Параметры запроса
            params = {
                "q": query,
                "apiKey": settings.NEWSAPI_KEY,
                "language": "en,ru",
                "sortBy": "publishedAt",
                "pageSize": settings.MAX_NEWS_PER_SOURCE,
            }

            async with httpx.AsyncClient(timeout=30.0) as client:
                response = await client.get(
                    f"{settings.NEWSAPI_URL}/everything", params=params
                )
                response.raise_for_status()
                data = response.json()

            news_items = []

            for article in data.get("articles", []):
                title = article.get("title", "")
                content = article.get("description", "") or article.get("content", "")
                url = article.get("url", "")
                author = article.get("author", "")
                image_url = article.get("urlToImage")

                # Дата публикации
                published_at = None
                if article.get("publishedAt"):
                    published_at = datetime.fromisoformat(
                        article["publishedAt"].replace("Z", "+00:00")
                    )

                # Проверка на дубликат
                content_hash = self.calculate_content_hash(title, content)
                if await self.check_duplicate(content_hash):
                    continue

                news_items.append(
                    {
                        "title": title,
                        "content": content,
                        "url": url,
                        "author": author,
                        "image_url": image_url,
                        "published_at": published_at,
                        "content_hash": content_hash,
                    }
                )

            logger.info(f"Collected {len(news_items)} news from NewsAPI")
            return news_items

        except Exception as e:
            logger.error(f"Error fetching from NewsAPI: {e}", exc_info=True)
            return []

    async def save_news_items(self, news_items: List[dict], source: NewsSource) -> int:
        """Сохранение новостей в БД"""
        saved_count = 0

        for item in news_items:
            try:
                news = NewsItem(
                    source_id=source.id,
                    title=item["title"],
                    content=item["content"],
                    url=item["url"],
                    author=item.get("author"),
                    image_url=item.get("image_url"),
                    category=source.category,
                    content_hash=item["content_hash"],
                    source_published_at=item.get("published_at"),
                    status=NewsStatus.PENDING,
                )

                self.db.add(news)
                saved_count += 1

            except Exception as e:
                logger.error(f"Error saving news item: {e}", item=item)

        await self.db.commit()
        logger.info(f"Saved {saved_count} news items to database")

        return saved_count

    async def collect_from_source(self, source: NewsSource) -> int:
        """Сбор новостей из одного источника"""
        logger.info(f"Starting collection from: {source.name}")

        # Парсинг в зависимости от типа источника
        if source.type == SourceType.RSS:
            news_items = await self.parse_rss_feed(source)
        elif source.type == SourceType.API:
            news_items = await self.fetch_from_newsapi(source)
        else:
            logger.warning(f"Unsupported source type: {source.type}")
            return 0

        # Сохранение в БД
        if news_items:
            saved_count = await self.save_news_items(news_items, source)

            # Обновление метаданных источника
            source.last_fetched_at = datetime.utcnow()
            source.last_success_at = datetime.utcnow()
            source.error_count = 0
            source.last_error = None
            await self.db.commit()

            return saved_count

        return 0

    async def collect_all(self, channel: Optional[NewsChannel] = None) -> dict:
        """Сбор новостей из всех активных источников"""
        logger.info(f"Starting collection for channel: {channel or 'all'}")

        # Получение активных источников
        query = select(NewsSource).where(NewsSource.is_active == True)
        if channel:
            query = query.where(NewsSource.category == channel)

        result = await self.db.execute(query)
        sources = result.scalars().all()

        logger.info(f"Found {len(sources)} active sources")

        # Сбор из каждого источника
        total_collected = 0
        results = {}

        for source in sources:
            try:
                count = await self.collect_from_source(source)
                total_collected += count
                results[source.name] = count
            except Exception as e:
                logger.error(f"Error collecting from {source.name}: {e}", exc_info=True)
                results[source.name] = 0

        logger.info(f"Collection complete: {total_collected} total news collected")

        return {
            "total_collected": total_collected,
            "sources": results,
            "timestamp": datetime.utcnow().isoformat(),
        }


async def initialize_default_sources(db: AsyncSession):
    """Инициализация источников новостей по умолчанию"""
    logger.info("Initializing default news sources")

    default_sources = [
        # Crypto sources
        {
            "name": "CoinTelegraph RSS",
            "type": SourceType.RSS,
            "url": "https://cointelegraph.com/rss",
            "category": NewsChannel.CRYPTO,
            "keywords": ["bitcoin", "ethereum", "crypto", "blockchain"],
        },
        {
            "name": "CoinDesk RSS",
            "type": SourceType.RSS,
            "url": "https://www.coindesk.com/arc/outboundfeeds/rss/",
            "category": NewsChannel.CRYPTO,
            "keywords": ["bitcoin", "crypto", "defi"],
        },
        {
            "name": "TechCrunch Crypto",
            "type": SourceType.RSS,
            "url": "https://techcrunch.com/tag/cryptocurrency/feed/",
            "category": NewsChannel.CRYPTO,
            "keywords": ["crypto", "blockchain", "web3"],
        },
        # Politics sources
        {
            "name": "RIA Novosti RSS",
            "type": SourceType.RSS,
            "url": "https://ria.ru/export/rss2/archive/index.xml",
            "category": NewsChannel.POLITICS,
            "keywords": ["russia", "kremlin", "putin"],
        },
        {
            "name": "TASS RSS",
            "type": SourceType.RSS,
            "url": "https://tass.ru/rss/v2.xml",
            "category": NewsChannel.POLITICS,
            "keywords": ["russia", "government", "politics"],
        },
        {
            "name": "BBC Russia",
            "type": SourceType.RSS,
            "url": "http://feeds.bbci.co.uk/russian/rss.xml",
            "category": NewsChannel.POLITICS,
            "keywords": ["russia", "ukraine", "moscow"],
        },
        # NewsAPI sources
        {
            "name": "NewsAPI Crypto",
            "type": SourceType.API,
            "url": settings.NEWSAPI_URL,
            "category": NewsChannel.CRYPTO,
            "keywords": ["crypto", "bitcoin"],
        },
        {
            "name": "NewsAPI Politics",
            "type": SourceType.API,
            "url": settings.NEWSAPI_URL,
            "category": NewsChannel.POLITICS,
            "keywords": ["russia", "kremlin"],
        },
    ]

    for source_data in default_sources:
        # Проверка существования
        result = await db.execute(
            select(NewsSource).where(NewsSource.name == source_data["name"])
        )
        existing = result.scalar_one_or_none()

        if not existing:
            source = NewsSource(**source_data)
            db.add(source)
            logger.info(f"Added source: {source_data['name']}")

    await db.commit()
    logger.info("Default sources initialized")
