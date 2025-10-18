from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession
from app.core.database import get_db
from pydantic import BaseModel

router = APIRouter()


class DashboardStats(BaseModel):
    total_news: int
    pending_news: int
    published_today: int
    ai_cost_today: float
    sources_active: int


@router.get("/dashboard", response_model=DashboardStats)
async def get_dashboard_stats(
    db: AsyncSession = Depends(get_db)
):
    """
    Получить статистику для dashboard админа
    """
    # TODO: Implement real stats from database
    return {
        "total_news": 5432,
        "pending_news": 23,
        "published_today": 45,
        "ai_cost_today": 12.34,
        "sources_active": 15
    }


@router.get("/logs")
async def get_system_logs(
    limit: int = 100,
    db: AsyncSession = Depends(get_db)
):
    """
    Получить системные логи
    """
    # TODO: Implement log retrieval
    return {
        "logs": [
            {
                "level": "info",
                "service": "collector",
                "message": "Collected 10 news from CoinDesk",
                "timestamp": "2025-01-18T12:00:00Z"
            }
        ]
    }
