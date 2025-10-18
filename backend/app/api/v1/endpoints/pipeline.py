from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.ext.asyncio import AsyncSession
from app.core.database import get_db
from typing import Optional

from app.services.collector import NewsCollector
from app.services.ai_analyzer import AIAnalyzer
from app.services.telegram_poster import TelegramPoster
from app.models import NewsChannel

router = APIRouter()


@router.post("/collect")
async def trigger_collection(
    channel: Optional[str] = None,
    db: AsyncSession = Depends(get_db)
):
    """
    Запустить сбор новостей вручную
    
    Args:
        channel: 'crypto' или 'politics' (опционально, по умолчанию все)
    """
    # Валидация канала
    news_channel = None
    if channel:
        try:
            news_channel = NewsChannel(channel)
        except ValueError:
            raise HTTPException(
                status_code=400, 
                detail=f"Invalid channel. Must be 'crypto' or 'politics'"
            )
    
    # Запуск сбора
    collector = NewsCollector(db)
    result = await collector.collect_all(channel=news_channel)
    
    return {
        "message": "News collection completed",
        "result": result
    }


@router.post("/analyze")
async def trigger_analysis(
    limit: int = 10,
    db: AsyncSession = Depends(get_db)
):
    """
    Запустить AI-анализ pending новостей вручную
    
    Args:
        limit: Максимальное количество новостей для анализа
    """
    analyzer = AIAnalyzer(db)
    result = await analyzer.analyze_pending_news(limit=limit)
    
    return {
        "message": "News analysis completed",
        "result": result
    }


@router.post("/post")
async def trigger_posting(
    limit: int = 5,
    db: AsyncSession = Depends(get_db)
):
    """
    Запустить публикацию analyzed новостей в Telegram вручную
    
    Args:
        limit: Максимальное количество новостей для публикации
    """
    poster = TelegramPoster(db)
    result = await poster.post_analyzed_news(limit=limit)
    
    return {
        "message": "News posting completed",
        "result": result
    }


@router.post("/pipeline")
async def run_full_pipeline(
    channel: Optional[str] = None,
    db: AsyncSession = Depends(get_db)
):
    """
    Запустить полный pipeline: collect → analyze → post
    
    Полезно для тестирования всей системы
    """
    # Валидация канала
    news_channel = None
    if channel:
        try:
            news_channel = NewsChannel(channel)
        except ValueError:
            raise HTTPException(
                status_code=400, 
                detail=f"Invalid channel. Must be 'crypto' or 'politics'"
            )
    
    results = {}
    
    # 1. Сбор
    collector = NewsCollector(db)
    results['collection'] = await collector.collect_all(channel=news_channel)
    
    # 2. Анализ
    analyzer = AIAnalyzer(db)
    results['analysis'] = await analyzer.analyze_pending_news(limit=20)
    
    # 3. Публикация
    poster = TelegramPoster(db)
    results['posting'] = await poster.post_analyzed_news(limit=10)
    
    return {
        "message": "Full pipeline completed",
        "results": results
    }
