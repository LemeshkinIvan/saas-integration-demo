package constants

import (
	"net/url"
	"time"
)

type CommonConfig struct {
	Postgres *PostgresConfig `yaml:"db"`
	Redis    *RedisConfig    `yaml:"cache"`
	Server   *ServerConfig   `yaml:"server"`
	Service  *ServiceConfig  `yaml:"service"`
}

type ServerConfig struct {
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type PostgresConfig struct {
	Password     string `yaml:"password"`
	UserName     string `yaml:"username"`
	Address      string `yaml:"address"`
	DatabaseName string `yaml:"name"`
}

func (c *PostgresConfig) GetConnectionString() string {
	login := url.QueryEscape(c.UserName) + ":" + url.QueryEscape(c.Password)
	res := "postgres://" + login + "@" + c.Address + "/" + c.DatabaseName
	return res
}

type RedisConfig struct {
	URL         string        `yaml:"url"`
	Password    string        `yaml:"password"`
	UserName    string        `yaml:"username"`
	DBNum       int           `yaml:"db_num"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

func (c *RedisConfig) GetConnectionString() string {
	return ""
}
