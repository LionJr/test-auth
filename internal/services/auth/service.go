package auth

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"test-auth/internal/config"
	"test-auth/pkg/token_manager"
)

type Service struct {
	config *config.AppConfig
	Logger *zap.Logger

	PgRepo      PgRepo
	RedisClient *redis.Client

	tokenManager *token_manager.TokenManager
}

func NewService(cfg *config.AppConfig, logger *zap.Logger, repo PgRepo, redisClient *redis.Client, tokenManager *token_manager.TokenManager) *Service {
	return &Service{
		config: cfg,
		Logger: logger,

		PgRepo:      repo,
		RedisClient: redisClient,

		tokenManager: tokenManager,
	}
}
