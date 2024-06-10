package config

import (
	"assignment/internal/errors"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

type ServiceConfig struct {
	URL           string
	Authenticator gin.HandlerFunc
	RateLimiter   func(c *gin.Context) error
	CacheDuration int
}

type Config struct {
	AppKey           string
	DBUser           string
	DBPassword       string
	DBHost           string
	DBPort           string
	DBName           string
	DefaultRateLimit int
	Services         map[string]*ServiceConfig
	REDIS_HOST       string
	REDIS_PORT       int
	REDIS_PASSWORD   string
	REDIs_DATABASE   int
}

var (
	globalConfig *Config
	once         sync.Once
)

func NewConfig() *errors.Error {
	if err := gotenv.Load(); err != nil {
		return errors.New(errors.ErrEnvNotFound)
	}

	defaultRateLimit, err := strconv.Atoi(os.Getenv("DEFAULT_RATE_LIMIT"))
	if err != nil {
		return errors.New(errors.ErrInvalidInput)
	}

	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {

		return errors.New(errors.ErrInvalidInput)
	}

	rediDatabase, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		return errors.New(errors.ErrInvalidInput)
	}

	once.Do(func() {
		globalConfig = &Config{
			AppKey:           os.Getenv("APP_KEY"),
			DBUser:           os.Getenv("DB_USER"),
			DBPassword:       os.Getenv("DB_PASSWORD"),
			DBHost:           os.Getenv("DB_HOST"),
			DBPort:           os.Getenv("DB_PORT"),
			DBName:           os.Getenv("DB_NAME"),
			DefaultRateLimit: defaultRateLimit,
			REDIS_HOST:       os.Getenv("REDIS_HOST"),
			REDIS_PORT:       redisPort,
			REDIS_PASSWORD:   os.Getenv("REDIS_PASSWORD"),
			REDIs_DATABASE:   rediDatabase,
		}
	})

	return nil
}

func GetConfig() *Config {
	return globalConfig
}
