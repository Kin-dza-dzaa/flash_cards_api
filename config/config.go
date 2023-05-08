package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Logger struct {
		// See https://pkg.go.dev/golang.org/x/exp/slog#Level.
		Level int `env:"LOGGER_LEVEL" env-default:"-4"`
	}

	OpenTelemetry struct {
		ServiceVersion string `env:"OTEL_SERVICE_VERSION" env-default:"0.3.4"`
		Environment    string `env:"OTEL_ENV" env-default:"development"`
		JaegerURL      string `env:"OTEL_JAEGER" env-default:"http://localhost:14268/api/traces"`
		ServiceName    string `env:"OTEL_SERVICE_NAME" env-default:"flash_cards_api"`
	}

	HTTP struct {
		Addr             string        `env:"HTTP_ADDR" env-default:"0.0.0.0:8000"`
		RateLimit        int           `env:"HTTP_RATE_LIMIT" env-default:"100"`
		RateWindow       time.Duration `env:"HTTP_RATE_WINDOW" env-default:"30s"`
		ReadTimeout      time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"5s"`
		WriteTimeout     time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"5s"`
		AllowCredentials bool          `env:"HTTP_ALLOW_CREDENTIALS" env-default:"true"`
		AllowedOrigins   []string      `env:"HTTP_ALLOWED_ORIGINS" env-separator:" " env-default:"http://localhost http://localhost:3000"`
		AllowedHeaders   []string      `env:"HTTP_ALLOWED_HEADERS" env-separator:" " env-default:"Content-Type Authorization"`
		AllowedMethods   []string      `env:"HTTP_ALLOWED_METHODS" env-separator:" " env-default:"POST GET PUT DELETE OPTIONS"`
		ShutdownTimeout  time.Duration `env:"HTTP_SHUT_DOWN_TIMEOUT" env-default:"10s"`
		// In seconds
		DefaultCorsDuration uint `env:"HTTP_DEFAULT_CORS_DURATION" env-default:"5"`
	}

	Postgres struct {
		URL         string `env-required:"true" env:"PG_URL" env-default:"postgresql://flash_cards:12345@localhost:5432/word_api"`
		MaxPoolSize int    `env:"PG_MAX_POOL_SIZE" env-default:"10"`
	}

	DictionaryAPI struct {
		URL             string `env:"GOOGLE_TRANSLATE_URL" env-default:"https://translate.google.com/_/TranslateWebserverUi/data/batchexecute"`
		DefaultSrcLang  string `env:"GOOGLE_TRANSLATE_DEFAULT_SRC" env-default:"en"`
		DefaultTrgtLang string `env:"GOOGLE_TRANSLATE_DEFAULT_TRGT" env-default:"ru"`
	}

	Cfg struct {
		OpenTelemetry OpenTelemetry
		GoogleAPI     DictionaryAPI
		PG            Postgres
		Logger        Logger
		HTTP          HTTP
	}
)

func ReadConfig() (Cfg, error) {
	cfg := Cfg{}
	err := cleanenv.ReadEnv(&cfg)
	return cfg, err
}
