package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Cfg struct {
		Logger
		HTTP
		Postgres
		DictionaryAPI
	}

	Logger struct {
		Level string `yaml:"log_level"`
	}

	HTTP struct {
		Host            string        `yaml:"host"`
		Port            string        `yaml:"port"`
		RateLimit       int           `yaml:"rate_limit"`
		RateWindow      time.Duration `yaml:"rate_window"`
		ReadTimeout     time.Duration `yaml:"read_timeout"`
		WriteTimeout    time.Duration `yaml:"write_timeout"`
		AllowedOrigins  []string      `yaml:"allowed_origins"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	}

	Postgres struct {
		PGURL       string `env-required:"true" env:"PG_URL"`
		MaxPoolSize int    `yaml:"max_pool_size"`
	}

	DictionaryAPI struct {
		DictionaryAPIURL string `yaml:"url"`
		DefaultSrcLang   string `yaml:"default_src_lang"`
		DefaultTrgtLang  string `yaml:"default_trgt_lang"`
	}
)

func ReadConfig(path string) (*Cfg, error) {
	cfg := new(Cfg)
	err := cleanenv.ReadConfig(path, cfg)
	return cfg, err
}
