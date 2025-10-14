package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	redis *redis.Client
}

func NewSessionRepository(redis *redis.Client) *SessionRepository {
	return &SessionRepository{redis: redis}
}

func (r *SessionRepository) SaveRefreshToken(ctx context.Context, userID, token string, expiry time.Duration) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, token)
	return r.redis.Set(ctx, key, "1", expiry).Err()
}

func (r *SessionRepository) ValidateRefreshToken(ctx context.Context, userID, token string) (bool, error) {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, token)
	result, err := r.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (r *SessionRepository) DeleteRefreshToken(ctx context.Context, userID, token string) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, token)
	return r.redis.Del(ctx, key).Err()
}

func (r *SessionRepository) DeleteAllUserTokens(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("refresh_token:%s:*", userID)
	iter := r.redis.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		if err := r.redis.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	return iter.Err()
}
