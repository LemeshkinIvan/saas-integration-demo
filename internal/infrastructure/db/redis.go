package db

import (
	"context"
	config "daos_core/internal/constants"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrKeyIsEmpty = errors.New("key is empty")

type RedisConstants struct {
	AccessToken         string
	AccessTokenTimeLive time.Duration
}

type RedisStorage struct {
	Instance *redis.Client
	Cfg      *config.RedisConfig
}

func CreateRedisInstance(cfg *config.RedisConfig) (*RedisStorage, error) {
	if err := validateInitialArgs(cfg); err != nil {
		return nil, err
	}

	fmt.Println(cfg.Password)

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.URL,
		Password: cfg.Password,
		DB:       cfg.DBNum,
		//Username:     cfg.UserName,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	return &RedisStorage{
		Instance: client,
		Cfg:      cfg,
	}, nil
}

func (r *RedisStorage) Ping() error {
	return r.Instance.Ping(context.Background()).Err()
}

func (r *RedisStorage) SetString(key string, value string, TTL time.Duration) error {
	if r.validateKey(key) {
		return ErrKeyIsEmpty
	}

	if err := r.Instance.Set(context.Background(), key, value, TTL).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisStorage) SetInt(key string, value int64, TTL time.Duration) error {
	if r.validateKey(key) {
		return ErrKeyIsEmpty
	}

	if err := r.Instance.Set(context.Background(), key, value, TTL).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisStorage) GetString(key string) (string, error) {
	if r.validateKey(key) {
		return "", ErrKeyIsEmpty
	}

	data, err := r.Instance.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}

	return data, nil
}

func (r *RedisStorage) GetInt(key string) (int64, error) {
	if r.validateKey(key) {
		return 0, ErrKeyIsEmpty
	}

	var data int64
	rawData, err := r.Instance.Get(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}

	data, err = strconv.ParseInt(rawData, 10, 64)
	if err != nil {
		return 0, err
	}

	return data, nil
}

func (r *RedisStorage) validateKey(value string) bool {
	return len(value) == 0
}

func (r *RedisStorage) Disconnect() error {
	return r.Instance.Close()
}

func validateInitialArgs(cfg *config.RedisConfig) error {
	if cfg == nil {
		return fmt.Errorf("")
	}

	if cfg.URL == "" {
		return fmt.Errorf("redis url is empty")
	}

	if cfg.Password == "" {
		return fmt.Errorf("redis pass is empty")
	}

	if cfg.DBNum < 0 {
		return fmt.Errorf("redis dbnum is empty")
	}

	// if cfg.UserName == "" {
	// 	return fmt.Errorf("redis user name is empty")
	// }

	if cfg.MaxRetries < 0 {
		return fmt.Errorf("redis max retries is empty")
	}
	if cfg.DialTimeout < 0 {
		return fmt.Errorf("redis dial timeout is empty")
	}

	if cfg.Timeout < 0 {
		return fmt.Errorf("redis timeout is empty")
	}
	return nil
}
