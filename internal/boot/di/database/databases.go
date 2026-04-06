package database

import (
	config "daos_core/internal/constants"
	"daos_core/internal/infrastructure/db"
	custom_logger "daos_core/internal/utils/logger"
	"errors"
)

type Container struct {
	PostgresDB *db.Postgres
	RedisDB    *db.RedisStorage
}

func NewDatabases(
	pEnv *config.PostgresConfig,
	r_env *config.RedisConfig,
) (*Container, error) {
	if pEnv == nil || r_env == nil {
		return nil, errors.New("database config or constant return nil")
	}

	postgres, err := db.Connect(pEnv)
	if err != nil {
		custom_logger.Logg.Error(err.Error())
		return nil, err
	}

	redis, err := db.CreateRedisInstance(r_env)
	if err != nil {
		return nil, err
	}

	// check connection
	if err := redis.Ping(); err != nil {
		return nil, err
	}

	return &Container{PostgresDB: postgres, RedisDB: redis}, nil
}
