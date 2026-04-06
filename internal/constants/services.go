package constants

import "time"

type ServiceConfig struct {
	Auth AuthConfig `yaml:"auth"`
}

type AuthConfig struct {
	AccessTTL  time.Duration `yaml:"access_ttl"`
	RefreshTTL time.Duration `yaml:"refresh_ttl"`
	Secret     []byte        `yaml:"secret"`
}
