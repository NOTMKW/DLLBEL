package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NOTMKW/DLLBEL/internal/models"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	client *redis.Client
	ctx context.Context
}

func NewRedisRepository(addr, password string, db int) *RedisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: db,
	})
	
	return &RedisRepository{
		client: rdb,
		ctx: 	context.Background(),
	}
}
func (r *RedisRepository) SaveRule(rule *models.Rule) error {
	data, err := json.Marshal(rule)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("rule:%s", rule.ID)
	return r.client.Set(r.ctx, key, data, 0).Err()
	}

func (r *RedisRepository) GetRule(id string) (*models.Rule, error) {
	key := fmt.Sprintf("rule:%s", id)
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var rule models.Rule
	if err := json.Unmarshal([]byte(data), &rule); err != nil {
		return nil, err
	}

	return &rule, nil
}

func (r *RedisRepository) GetAllRules() ([]*models.Rule, error) {
	keys, err := r.client.Keys(r.ctx, "rule:*").Result()
	if err != nil {
		return nil, err
	}

	rules := make([]*models.Rule, 0, len(keys))
	for _, key := range keys {
		data, err := r.client.Get(r.ctx, key).Result()
		if err != nil {
			continue
		}

		var rule models.Rule
		if err := json.Unmarshal([]byte(data), &rule); err != nil {
			continue
		}
		rules = append(rules, &rule)
	}

	return rules, nil
}

func (r *RedisRepository) DeleteRule(id string) error {
	key := fmt.Sprintf("rule:%s", id)
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisRepository) SaveUserState(state *models.UserState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("user_state:%s", state.UserID)
	return r.client.Set(r.ctx, key, data, 0).Err()
}

func (r *RedisRepository) GetUserState(userID string) (*models.UserState, error) {
	key := fmt.Sprintf("user_state:%s", userID)
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var state models.UserState
	if err := json.Unmarshal([]byte(data), &state); err != nil {
		return nil, err
	}

	return &state, nil
}

func (r *RedisRepository) Close() error {
	return r.client.Close()
}