package session

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisInterface interface {
	Get(ctx context.Context, sessionID string) (map[string]interface{}, error)
	Set(ctx context.Context, sessionID string, data map[string]interface{}, ttl time.Duration) error
	Delete(ctx context.Context, sessionID string) error
}

type RedisSessionStore struct {
	client *redis.Client
}

func NewRedisSessionStore(client *redis.Client) *RedisSessionStore {
	return &RedisSessionStore{client: client}
}

func (r *RedisSessionStore) Get(ctx context.Context, sessionID string) (map[string]interface{}, error) {
	data, err := r.client.Get(ctx, sessionID).Result()
	if errors.Is(err, redis.Nil) {
		return nil, errors.New("session not found")
	}
	if err != nil {
		return nil, err
	}

	var sessionData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &sessionData); err != nil {
		return nil, err
	}

	return sessionData, nil
}

func (r *RedisSessionStore) Set(ctx context.Context, sessionID string, data map[string]interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, sessionID, jsonData, ttl).Err()
}

func (r *RedisSessionStore) Delete(ctx context.Context, sessionID string) error {
	return r.client.Del(ctx, sessionID).Err()
}
