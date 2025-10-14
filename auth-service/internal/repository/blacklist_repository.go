package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type BlacklistRepository struct {
	redis *redis.Client
}

func NewBlacklistRepository(redis *redis.Client) *BlacklistRepository {
	return &BlacklistRepository{redis: redis}
}

// AddToBlacklist добавляет токен в черный список
func (r *BlacklistRepository) AddToBlacklist(ctx context.Context, token string, expiry time.Duration) error {
	key := fmt.Sprintf("blacklist:token:%s", token)
	return r.redis.Set(ctx, key, "1", expiry).Err()
}

// IsBlacklisted проверяет, находится ли токен в черном списке
func (r *BlacklistRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:token:%s", token)
	result, err := r.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// RemoveFromBlacklist удаляет токен из черного списка
func (r *BlacklistRepository) RemoveFromBlacklist(ctx context.Context, token string) error {
	key := fmt.Sprintf("blacklist:token:%s", token)
	return r.redis.Del(ctx, key).Err()
}
