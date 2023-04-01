package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Logger struct {
		Level string `env:"LOGGER_LEVEL" env-default:"debug"`
	}

	HTTP struct {
		Addr            string        `env:"HTTP_ADDR" env-default:"0.0.0.0:8000"`
		RateLimit       int           `env:"HTTP_RATE_LIMIT" env-default:"100"`
		RateWindow      time.Duration `env:"HTTP_RATE_WINDOW" env-default:"30s"`
		ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"5s"`
		WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"5s"`
		AllowedOrigins  []string      `env:"HTTP_ALLOWED_ORIGINS" env-separator:" " env-default:"http://localhost http://localhost:3000"`
		AllowedHeaders  []string      `env:"HTTP_ALLOWED_HEADERS" env-separator:" " env-default:"Content-Type Authorization"`
		ShutdownTimeout time.Duration `env:"HTTP_SHUT_DOWN_TIMEOUT" env-default:"10s"`
	}

	Postgres struct {
		URL         string `env-required:"true" env:"PG_URL"`
		MaxPoolSize int    `env:"PG_MAX_POOL_SIZE" env-default:"10"`
	}

	DictionaryAPI struct {
		URL             string `env:"GOOGLE_TRANSLATE_URL" env-default:"https://translate.google.com/_/TranslateWebserverUi/data/batchexecute"`
		DefaultSrcLang  string `env:"GOOGLE_TRANSLATE_DEFAULT_SRC" env-default:"en"`
		DefaultTrgtLang string `env:"GOOGLE_TRANSLATE_DEFAULT_TRGT" env-default:"ru"`
	}

	Cfg struct {
		Logger    Logger
		HTTP      HTTP
		PG        Postgres
		GoogleApi DictionaryAPI
	}
)

func ReadConfig() (Cfg, error) {
	cfg := Cfg{}
	err := cleanenv.ReadEnv(&cfg)
	return cfg, err
}
