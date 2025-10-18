from datetime import datetime
from typing import List, Optional
from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query
from pydantic import BaseModel
from sqlalchemy import func, or_, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.database import get_db
from app.models import NewsChannel, NewsItem, NewsStatus
from app.services.telegram_poster import approve_news_for_posting, reject_news

router = APIRouter()


class NewsResponse(BaseModel):
    id: str
    title: str
    summary: Optional[str]
    ai_summary: Optional[str]
    ai_insights: Optional[str]
    url: str
    image_url: Optional[str]
    category: str
    relevance_score: Optional[float]
    published_at: Optional[datetime]
    status: str
    hashtags: Optional[List[str]]

    class Config:
        from_attributes = True


class NewsListResponse(BaseModel):
    items: List[NewsResponse]
    total: int
    page: int
    per_page: int


@router.get("/", response_model=NewsListResponse)
async def get_news(
    page: int = Query(1, ge=1),
    per_page: int = Query(20, ge=1, le=100),
    category: Optional[str] = None,
    status: Optional[str] = None,
    search: Optional[str] = None,
    db: AsyncSession = Depends(get_db),
):
    """
    Получить список новостей с пагинацией и фильтрами
    """
    # Базовый запрос
    query = select(NewsItem)

    # Фильтры
    filters = []
    if category:
        try:
            filters.append(NewsItem.category == NewsChannel(category))
        except ValueError:
            raise HTTPException(status_code=400, detail="Invalid category")

    if status:
        try:
            filters.append(NewsItem.status == NewsStatus(status))
        except ValueError:
            raise HTTPException(status_code=400, detail="Invalid status")

    if search:
        # Поиск по заголовку или содержимому
        search_filter = or_(
            NewsItem.title.ilike(f"%{search}%"), NewsItem.content.ilike(f"%{search}%")
        )
        filters.append(search_filter)

    if filters:
        query = query.where(*filters)

    # Подсчёт общего количества
    count_query = select(func.count()).select_from(NewsItem)
    if filters:
        count_query = count_query.where(*filters)

    total_result = await db.execute(count_query)
    total = total_result.scalar()

    # Пагинация и сортировка
    query = query.order_by(NewsItem.created_at.desc())
    query = query.offset((page - 1) * per_page).limit(per_page)

    result = await db.execute(query)
    news_items = result.scalars().all()

    # Конвертация в response модель
    items = [
        NewsResponse(
            id=str(item.id),
            title=item.title,
            summary=(
                item.content[:200] + "..." if len(item.content) > 200 else item.content
            ),
            ai_summary=item.ai_summary,
            ai_insights=item.ai_insights,
            url=item.url,
            image_url=item.image_url,
            category=item.category.value,
            relevance_score=item.relevance_score,
            published_at=item.published_at,
            status=item.status.value,
            hashtags=item.ai_hashtags,
        )
        for item in news_items
    ]

    return {"items": items, "total": total, "page": page, "per_page": per_page}


@router.get("/{news_id}", response_model=NewsResponse)
async def get_news_by_id(news_id: str, db: AsyncSession = Depends(get_db)):
    """
    Получить детали конкретной новости
    """
    try:
        news_uuid = UUID(news_id)
    except ValueError:
        raise HTTPException(status_code=400, detail="Invalid news ID format")

    result = await db.execute(select(NewsItem).where(NewsItem.id == news_uuid))
    news = result.scalar_one_or_none()

    if not news:
        raise HTTPException(status_code=404, detail="News not found")

    return NewsResponse(
        id=str(news.id),
        title=news.title,
        summary=news.content,
        ai_summary=news.ai_summary,
        ai_insights=news.ai_insights,
        url=news.url,
        image_url=news.image_url,
        category=news.category.value,
        relevance_score=news.relevance_score,
        published_at=news.published_at,
        status=news.status.value,
        hashtags=news.ai_hashtags,
    )


@router.post("/{news_id}/approve")
async def approve_news(news_id: str, db: AsyncSession = Depends(get_db)):
    """
    Одобрить новость для публикации
    """
    success = await approve_news_for_posting(news_id, db)

    if not success:
        raise HTTPException(status_code=400, detail="Failed to approve news")

    return {"message": f"News {news_id} approved and posted to Telegram"}


@router.post("/{news_id}/reject")
async def reject_news_endpoint(news_id: str, db: AsyncSession = Depends(get_db)):
    """
    Отклонить новость
    """
    success = await reject_news(news_id, db)

    if not success:
        raise HTTPException(status_code=400, detail="Failed to reject news")

    return {"message": f"News {news_id} rejected"}


@router.delete("/{news_id}")
async def delete_news(news_id: str, db: AsyncSession = Depends(get_db)):
    """
    Удалить новость
    """
    try:
        news_uuid = UUID(news_id)
    except ValueError:
        raise HTTPException(status_code=400, detail="Invalid news ID format")

    result = await db.execute(select(NewsItem).where(NewsItem.id == news_uuid))
    news = result.scalar_one_or_none()

    if not news:
        raise HTTPException(status_code=404, detail="News not found")

    await db.delete(news)
    await db.commit()

    return {"message": f"News {news_id} deleted"}
