from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession
from app.core.database import get_db
from pydantic import BaseModel, HttpUrl
from typing import List

router = APIRouter()


class SourceCreate(BaseModel):
    name: str
    type: str  # 'rss', 'api', 'scraping'
    url: HttpUrl
    category: str  # 'crypto', 'politics', 'it', 'world'


class SourceResponse(BaseModel):
    id: int
    name: str
    type: str
    url: str
    category: str
    is_active: bool


@router.get("/", response_model=List[SourceResponse])
async def get_sources(
    db: AsyncSession = Depends(get_db)
):
    """
    Получить список всех источников новостей
    """
    # TODO: Implement database query
    return [
        {
            "id": 1,
            "name": "CoinDesk RSS",
            "type": "rss",
            "url": "https://www.coindesk.com/arc/outboundfeeds/rss/",
            "category": "crypto",
            "is_active": True
        }
    ]


@router.post("/", response_model=SourceResponse)
async def create_source(
    source: SourceCreate,
    db: AsyncSession = Depends(get_db)
):
    """
    Добавить новый источник новостей
    """
    # TODO: Implement source creation
    return {
        "id": 1,
        "name": source.name,
        "type": source.type,
        "url": str(source.url),
        "category": source.category,
        "is_active": True
    }


@router.post("/collect")
async def trigger_collection(
    db: AsyncSession = Depends(get_db)
):
    """
    Запустить сбор новостей вручную
    """
    # TODO: Trigger Celery task for news collection
    return {"message": "News collection started", "task_id": "abc123"}
