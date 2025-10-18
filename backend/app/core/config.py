from typing import List

from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    """Application settings"""

    # App
    APP_NAME: str = "NewsHub AI"
    ENVIRONMENT: str = "development"
    DEBUG: bool = False

    # Database
    DATABASE_URL: str

    # Redis
    REDIS_URL: str = "redis://localhost:6379/0"

    # RabbitMQ
    RABBITMQ_URL: str = "amqp://guest:guest@localhost:5672/"

    # JWT
    JWT_SECRET_KEY: str
    JWT_ALGORITHM: str = "HS256"
    JWT_ACCESS_TOKEN_EXPIRE_MINUTES: int = 15
    JWT_REFRESH_TOKEN_EXPIRE_DAYS: int = 7

    # OpenRouter
    OPENROUTER_API_KEY: str
    OPENROUTER_API_URL: str = "https://openrouter.ai/api/v1"
    OPENROUTER_MODEL: str = "openai/gpt-4-turbo-preview"

    # Freepik
    FREEPIK_API_KEY: str
    FREEPIK_API_URL: str = "https://api.freepik.com/v1"

    # NewsAPI
    NEWSAPI_KEY: str
    NEWSAPI_URL: str = "https://newsapi.org/v2"

    # Telegram
    TELEGRAM_BOT_TOKEN: str
    TELEGRAM_CRYPTO_CHANNEL: str = "@crypto_ainews"
    TELEGRAM_POLITICS_CHANNEL: str = "@kremlin_digest"
    TELEGRAM_ADMIN_CHAT_ID: str

    # CORS
    ALLOWED_ORIGINS: List[str] = ["http://localhost:3000", "http://localhost"]

    # Rate Limiting
    RATE_LIMIT_PER_MINUTE: int = 60

    # Collection Settings
    COLLECT_INTERVAL_MINUTES: int = 15
    MAX_NEWS_PER_SOURCE: int = 50

    # AI Settings
    AI_TIMEOUT_SECONDS: int = 30
    AI_MAX_RETRIES: int = 3
    AI_IMPORTANCE_THRESHOLD: int = 4

    class Config:
        env_file = ".env"
        case_sensitive = True


settings = Settings()
