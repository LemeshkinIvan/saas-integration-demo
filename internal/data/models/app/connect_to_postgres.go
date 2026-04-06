package app_models

import "time"

type PostgresConnectionDTO struct {
	Password     string
	Name         string
	Port         uint
	Url          string
	DatabaseName string
}

type RedisConfig struct {
	Address     string
	Password    string
	User        string
	DB          int
	MaxRetries  int
	DialTimeout time.Duration
	Timeout     time.Duration
}
