from fastapi import APIRouter, Depends
from pydantic import BaseModel
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.database import get_db

router = APIRouter()


class TelegramPostRequest(BaseModel):
    news_id: str
    channel: str  # '@crypto_ainews' or '@kremlin_digest'


@router.post("/post")
async def post_to_telegram(
    request: TelegramPostRequest, db: AsyncSession = Depends(get_db)
):
    """
    Опубликовать новость в Telegram-канал
    """
    # TODO: Implement Telegram posting via Celery task
    return {"message": f"News posted to {request.channel}", "message_id": 12345}


@router.get("/stats")
async def get_telegram_stats(db: AsyncSession = Depends(get_db)):
    """
    Получить статистику по Telegram-каналам
    """
    # TODO: Implement stats collection
    return {
        "crypto_ainews": {"posts_today": 15, "total_posts": 1234, "avg_views": 5678},
        "kremlin_digest": {"posts_today": 12, "total_posts": 987, "avg_views": 3456},
    }
