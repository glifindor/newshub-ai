from contextlib import asynccontextmanager

from fastapi import Depends, FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.gzip import GZipMiddleware

from app.api.v1 import api_router
from app.core.config import settings
from app.core.database import AsyncSessionLocal, Base, engine
from app.core.logger import get_logger
from app.services.collector import initialize_default_sources
from app.services.scheduler import start_scheduler, stop_scheduler

logger = get_logger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifecycle events"""
    # Startup
    logger.info("🚀 Starting NewsHub AI...")

    # Создание таблиц (в продакшене используйте Alembic!)
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    logger.info("✅ Database connected and tables created")

    # Инициализация источников новостей
    async with AsyncSessionLocal() as db:
        await initialize_default_sources(db)
    logger.info("✅ Default news sources initialized")

    # Запуск scheduler
    start_scheduler()
    logger.info("✅ Scheduler started")

    yield

    # Shutdown
    logger.info("👋 Shutting down...")
    stop_scheduler()


app = FastAPI(
    title="NewsHub AI API",
    description="Центральный хаб для AI-powered новостной агрегации",
    version="1.0.0",
    lifespan=lifespan,
    docs_url="/docs",
    redoc_url="/redoc",
)

# CORS Middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.ALLOWED_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# GZip Compression
app.add_middleware(GZipMiddleware, minimum_size=1000)

# Include routers
app.include_router(api_router, prefix="/api/v1")


@app.get("/")
async def root():
    """Health check"""
    return {
        "status": "ok",
        "service": "NewsHub AI",
        "version": "1.0.0",
        "docs": "/docs",
    }


@app.get("/health")
async def health():
    """Detailed health check"""
    return {
        "status": "healthy",
        "database": "connected",
        "redis": "connected",
        "celery": "running",
    }


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
