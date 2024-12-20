package session

import (
	"2024_2_FIGHT-CLUB/domain"
	"context"
	"errors"
	"github.com/mailru/easyjson"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisInterface interface {
	Get(ctx context.Context, sessionID string) (domain.SessionData, error)
	Set(ctx context.Context, sessionID string, data domain.SessionData, ttl time.Duration) error
	Delete(ctx context.Context, sessionID string) error
}

type RedisSessionStore struct {
	client *redis.Client
}

func NewRedisSessionStore(client *redis.Client) *RedisSessionStore {
	return &RedisSessionStore{client: client}
}

func (r *RedisSessionStore) Get(ctx context.Context, sessionID string) (domain.SessionData, error) {
	data, err := r.client.Get(ctx, sessionID).Result()
	if errors.Is(err, redis.Nil) {
		return domain.SessionData{}, errors.New("session not found")
	}
	if err != nil {
		return domain.SessionData{}, err
	}

	var sessionData domain.SessionData
	if err := easyjson.Unmarshal([]byte(data), &sessionData); err != nil {
		return domain.SessionData{}, err
	}

	return sessionData, nil
}

func (r *RedisSessionStore) Set(ctx context.Context, sessionID string, data domain.SessionData, ttl time.Duration) error {
	jsonData, err := easyjson.Marshal(data)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, sessionID, jsonData, ttl).Err()
}

func (r *RedisSessionStore) Delete(ctx context.Context, sessionID string) error {
	return r.client.Del(ctx, sessionID).Err()
}
