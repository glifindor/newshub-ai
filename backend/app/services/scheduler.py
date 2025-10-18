"""
Scheduler для автоматического сбора и обработки новостей
"""
import asyncio
from datetime import datetime
from apscheduler.schedulers.asyncio import AsyncIOScheduler
from apscheduler.triggers.interval import IntervalTrigger

from app.core.database import AsyncSessionLocal
from app.services.collector import NewsCollector
from app.services.ai_analyzer import AIAnalyzer
from app.services.telegram_poster import TelegramPoster
from app.core.config import settings
from app.core.logger import get_logger

logger = get_logger(__name__)

# Глобальный scheduler
scheduler = AsyncIOScheduler()


async def collect_news_job():
    """Задача сбора новостей"""
    logger.info("=== Starting news collection job ===")
    
    async with AsyncSessionLocal() as db:
        try:
            collector = NewsCollector(db)
            result = await collector.collect_all()
            
            logger.info(
                "News collection job completed",
                total_collected=result['total_collected'],
                sources=result['sources']
            )
            
        except Exception as e:
            logger.error(f"Error in news collection job: {e}", exc_info=True)


async def analyze_news_job():
    """Задача AI-анализа новостей"""
    logger.info("=== Starting news analysis job ===")
    
    async with AsyncSessionLocal() as db:
        try:
            analyzer = AIAnalyzer(db)
            result = await analyzer.analyze_pending_news(limit=20)
            
            logger.info(
                "News analysis job completed",
                total=result['total'],
                analyzed=result['analyzed'],
                rejected=result['rejected'],
                failed=result['failed']
            )
            
        except Exception as e:
            logger.error(f"Error in news analysis job: {e}", exc_info=True)


async def post_news_job():
    """Задача публикации новостей в Telegram"""
    logger.info("=== Starting news posting job ===")
    
    async with AsyncSessionLocal() as db:
        try:
            poster = TelegramPoster(db)
            
            # Сначала обрабатываем модерацию
            moderation_result = await poster.handle_moderation_requests()
            logger.info(
                "Moderation requests handled",
                total=moderation_result['total'],
                notified=moderation_result['notified']
            )
            
            # Затем постим проанализированные
            post_result = await poster.post_analyzed_news(limit=10)
            logger.info(
                "News posting job completed",
                total=post_result['total'],
                posted=post_result['posted'],
                failed=post_result['failed']
            )
            
        except Exception as e:
            logger.error(f"Error in news posting job: {e}", exc_info=True)


def start_scheduler():
    """Запуск scheduler"""
    logger.info("Starting news scheduler")
    
    # Сбор новостей каждые 10 минут
    scheduler.add_job(
        collect_news_job,
        trigger=IntervalTrigger(minutes=settings.COLLECT_INTERVAL_MINUTES),
        id='collect_news',
        name='Collect news from sources',
        replace_existing=True,
    )
    
    # AI-анализ каждые 5 минут
    scheduler.add_job(
        analyze_news_job,
        trigger=IntervalTrigger(minutes=5),
        id='analyze_news',
        name='Analyze news with AI',
        replace_existing=True,
    )
    
    # Публикация каждые 7 минут
    scheduler.add_job(
        post_news_job,
        trigger=IntervalTrigger(minutes=7),
        id='post_news',
        name='Post news to Telegram',
        replace_existing=True,
    )
    
    # Запуск scheduler
    scheduler.start()
    logger.info("Scheduler started successfully")


def stop_scheduler():
    """Остановка scheduler"""
    logger.info("Stopping scheduler")
    scheduler.shutdown()
