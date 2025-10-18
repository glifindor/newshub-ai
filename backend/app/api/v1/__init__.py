from fastapi import APIRouter

from app.api.v1.endpoints import admin, auth, news, pipeline, sources, telegram

api_router = APIRouter()

# Include all endpoint routers
api_router.include_router(auth.router, prefix="/auth", tags=["Authentication"])
api_router.include_router(news.router, prefix="/news", tags=["News"])
api_router.include_router(sources.router, prefix="/sources", tags=["Sources"])
api_router.include_router(telegram.router, prefix="/telegram", tags=["Telegram"])
api_router.include_router(admin.router, prefix="/admin", tags=["Admin"])
api_router.include_router(pipeline.router, prefix="/pipeline", tags=["Pipeline"])
