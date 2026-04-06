package utils

import (
	"daos_core/internal/constants"
	"daos_core/internal/utils/jwt"
	"time"
)

type Container struct {
	JWT jwt.JWTUtil
}

func RegisterAll(cfg *constants.AuthConfig) (*Container, error) {
	jwt := jwt.NewJWTUtil(
		map[string]time.Duration{
			"access":  cfg.AccessTTL,
			"refresh": cfg.RefreshTTL,
		},
		cfg.Secret,
	)

	return &Container{JWT: jwt}, nil
}
