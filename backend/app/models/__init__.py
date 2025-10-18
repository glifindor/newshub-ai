"""
Database models для NewsHub AI
"""

import enum
import json
import uuid
from datetime import datetime

from sqlalchemy import Boolean, Column, DateTime
from sqlalchemy import Enum as SQLEnum
from sqlalchemy import Float, Integer, String, Text
from sqlalchemy.dialects.postgresql import ARRAY, UUID
from sqlalchemy.types import Text as TextType
from sqlalchemy.types import TypeDecorator

from app.core.database import Base

# Очищаем metadata при модульном импорте, чтобы избежать дублирования таблиц в тестовом окружении
Base.metadata.clear()


class StringArray(TypeDecorator):
    """Cross-dialect friendly list of strings stored as JSON when ARRAY is unavailable."""

    cache_ok = True
    impl = TextType

    def load_dialect_impl(self, dialect):
        if dialect.name == "postgresql":
            return dialect.type_descriptor(ARRAY(String))
        return dialect.type_descriptor(TextType)

    def process_bind_param(self, value, dialect):
        if value is None:
            return None
        if dialect.name == "postgresql":
            return value
        if isinstance(value, str):
            return value
        return json.dumps(value)

    def process_result_value(self, value, dialect):
        if value is None:
            return None
        if dialect.name == "postgresql":
            return value
        try:
            return json.loads(value)
        except (TypeError, json.JSONDecodeError):
            return value


class GUID(TypeDecorator):
    """Platform-independent GUID type."""

    cache_ok = True
    impl = UUID

    def load_dialect_impl(self, dialect):
        if dialect.name == "postgresql":
            return dialect.type_descriptor(UUID(as_uuid=True))
        return dialect.type_descriptor(String(36))

    def process_bind_param(self, value, dialect):
        if value is None:
            return value
        if dialect.name == "postgresql":
            return value
        if isinstance(value, uuid.UUID):
            return str(value)
        return str(uuid.UUID(str(value)))

    def process_result_value(self, value, dialect):
        if value is None:
            return value
        if dialect.name == "postgresql":
            return value
        return uuid.UUID(value)


class NewsStatus(str, enum.Enum):
    """Статусы новости"""

    PENDING = "pending"  # Собрана, ждёт AI-анализа
    ANALYZED = "analyzed"  # AI-анализ завершён
    APPROVED = "approved"  # Одобрена админом
    PUBLISHED = "published"  # Опубликована в Telegram
    REJECTED = "rejected"  # Отклонена (низкий рейтинг)


class NewsChannel(str, enum.Enum):
    """Каналы публикации"""

    CRYPTO = "crypto"  # @crypto_ainews
    POLITICS = "politics"  # @kremlin_digest


class SourceType(str, enum.Enum):
    """Типы источников"""

    RSS = "rss"
    API = "api"
    SCRAPING = "scraping"


class NewsSource(Base):
    """Источники новостей"""

    __tablename__ = "news_sources"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String(255), nullable=False)
    type = Column(SQLEnum(SourceType), nullable=False)
    url = Column(Text, nullable=False)
    category = Column(SQLEnum(NewsChannel), nullable=False)
    is_active = Column(Boolean, default=True)

    # Keywords для фильтрации
    keywords = Column(StringArray(), nullable=True)

    # Метаданные
    last_fetched_at = Column(DateTime, nullable=True)
    last_success_at = Column(DateTime, nullable=True)
    error_count = Column(Integer, default=0)
    last_error = Column(Text, nullable=True)

    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)


class NewsItem(Base):
    """Новости"""

    __tablename__ = "news_items"

    id = Column(GUID(), primary_key=True, default=uuid.uuid4, index=True)
    source_id = Column(Integer, nullable=False)

    # Контент
    title = Column(String(500), nullable=False)
    content = Column(Text, nullable=False)
    url = Column(Text, nullable=False, unique=True)
    image_url = Column(Text, nullable=True)
    author = Column(String(255), nullable=True)

    # Метаданные
    category = Column(SQLEnum(NewsChannel), nullable=False)
    tags = Column(StringArray(), nullable=True)
    language = Column(String(10), default="ru")

    # AI-анализ
    ai_summary = Column(Text, nullable=True)  # Краткий тизер
    ai_insights = Column(Text, nullable=True)  # Инсайты для инвесторов/политики
    ai_hashtags = Column(StringArray(), nullable=True)  # Хэштеги
    relevance_score = Column(Float, nullable=True)  # Рейтинг релевантности (0-10)

    # Статус
    status = Column(SQLEnum(NewsStatus), default=NewsStatus.PENDING, index=True)
    requires_moderation = Column(Boolean, default=False)

    # Публикация
    published_at = Column(DateTime, nullable=True)
    telegram_message_id = Column(Integer, nullable=True)
    telegram_channel = Column(String(100), nullable=True)

    # Дедупликация
    content_hash = Column(String(64), unique=True, index=True)

    # Timestamps
    source_published_at = Column(DateTime, nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow, index=True)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)


class User(Base):
    """Пользователи/Админы"""

    __tablename__ = "users"

    id = Column(Integer, primary_key=True, index=True)
    username = Column(String(100), unique=True, nullable=False)
    email = Column(String(255), unique=True, nullable=False)
    hashed_password = Column(String(255), nullable=False)
    role = Column(String(50), default="moderator")  # admin, moderator, viewer
    is_active = Column(Boolean, default=True)
    telegram_chat_id = Column(Integer, nullable=True)

    created_at = Column(DateTime, default=datetime.utcnow)


class AITask(Base):
    """Задачи AI-обработки"""

    __tablename__ = "ai_tasks"

    id = Column(GUID(), primary_key=True, default=uuid.uuid4)
    news_item_id = Column(GUID(), nullable=False, index=True)

    task_type = Column(String(50), nullable=False)  # analyze, generate_image, summarize
    status = Column(
        String(50), default="pending"
    )  # pending, processing, completed, failed

    provider = Column(String(50), nullable=True)  # openrouter, freepik
    model = Column(String(100), nullable=True)

    # Результаты
    output_data = Column(Text, nullable=True)
    error_message = Column(Text, nullable=True)

    # Метрики
    processing_time_ms = Column(Integer, nullable=True)
    cost_usd = Column(Float, nullable=True)

    created_at = Column(DateTime, default=datetime.utcnow)
    completed_at = Column(DateTime, nullable=True)


class SystemLog(Base):
    """Системные логи"""

    __tablename__ = "system_logs"

    id = Column(Integer, primary_key=True, index=True)
    level = Column(String(20), nullable=False)  # info, warning, error, critical
    service = Column(
        String(100), nullable=False
    )  # collector, ai_analyzer, telegram_poster
    message = Column(Text, nullable=False)

    created_at = Column(DateTime, default=datetime.utcnow, index=True)
