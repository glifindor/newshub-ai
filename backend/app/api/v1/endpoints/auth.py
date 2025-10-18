from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from pydantic import BaseModel, EmailStr
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.database import get_db
from app.core.security import create_access_token, create_refresh_token, verify_password

router = APIRouter()
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="api/v1/auth/login")


class UserLogin(BaseModel):
    username: str
    password: str


class Token(BaseModel):
    access_token: str
    refresh_token: str
    token_type: str = "bearer"


class UserResponse(BaseModel):
    id: int
    username: str
    email: str
    role: str


@router.post("/login", response_model=Token)
async def login(
    form_data: OAuth2PasswordRequestForm = Depends(), db: AsyncSession = Depends(get_db)
):
    """
    Login endpoint - получение JWT токенов
    """
    # TODO: Implement actual database user lookup
    # Временная заглушка для демонстрации
    if form_data.username == "admin" and form_data.password == "admin":
        access_token = create_access_token(
            data={"sub": form_data.username, "role": "admin"}
        )
        refresh_token = create_refresh_token(data={"sub": form_data.username})

        return {
            "access_token": access_token,
            "refresh_token": refresh_token,
            "token_type": "bearer",
        }

    raise HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Incorrect username or password",
        headers={"WWW-Authenticate": "Bearer"},
    )


@router.post("/refresh", response_model=Token)
async def refresh_token(refresh_token: str):
    """
    Refresh access token using refresh token
    """
    # TODO: Implement token refresh logic
    pass


@router.get("/me", response_model=UserResponse)
async def get_current_user(
    token: str = Depends(oauth2_scheme), db: AsyncSession = Depends(get_db)
):
    """
    Get current authenticated user
    """
    # TODO: Implement user extraction from token
    return {"id": 1, "username": "admin", "email": "admin@newshub.ai", "role": "admin"}


@router.post("/logout")
async def logout(token: str = Depends(oauth2_scheme)):
    """
    Logout - invalidate token
    """
    # TODO: Implement token blacklisting with Redis
    return {"message": "Successfully logged out"}
