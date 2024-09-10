package main

import (
	"context"
	"flag"
	"go.uber.org/zap"
	"log"
	"test-auth/internal/app"
	"test-auth/internal/config"
	"test-auth/internal/env"
	"test-auth/internal/repositories/postgres"
	"test-auth/internal/repositories/redis"
	"test-auth/internal/services/auth"
	"test-auth/internal/smtp"
	"test-auth/pkg/token_manager"
)

const (
	Name        = "test-auth"
	CurrentPath = "."
)

// @title           Test task Back-Dev Auth API
// @version         1.0
// @description     API server for Auth API

// @host      192.168.77.110:8080
// @BasePath  /api

func main() {
	ctx := context.Background()

	var version, environment, logLevel string
	flag.StringVar(&version, "v", "", "version")
	flag.StringVar(&environment, "e", "", "environment")
	flag.StringVar(&logLevel, "ll", "info", "logging level")
	flag.Parse()

	ctx = context.WithValue(ctx, env.Name, Name)
	ctx = context.WithValue(ctx, env.Version, version)
	ctx = context.WithValue(ctx, env.Environment, environment)

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("failed to init logger", err.Error())
	}

	cfg, err := config.NewAppConfig(CurrentPath + "/configs/" + environment + ".yml")
	if err != nil {
		logger.Fatal("failed to read config: ", zap.Error(err))
	}

	pgClient, err := postgres.NewPostgresDB(cfg.Postgres)
	if err != nil {
		logger.Fatal("error while connecting to Postgres db: ", zap.Error(err))
	}

	redisClient, err := redis.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		logger.Fatal("error while connecting to Redis", zap.Error(err))
	}

	tokenManager, err := token_manager.NewManager(cfg.Token.AccessSecretKey, cfg.Token.RefreshSecretKey)
	if err != nil {
		logger.Fatal("error while initializing token manager", zap.Error(err))
	}

	logger.Info(
		"flags",
		zap.String("version", version),
		zap.String("environment", environment),
		zap.String("log_level", logLevel),
	)

	newSmtp := smtp.NewSmtp(cfg.SMTP, redisClient, tokenManager)

	application := &app.Application{
		Config:         cfg,
		Logger:         logger,
		PostgresClient: pgClient,
		RedisClient:    redisClient,
		TokenManager:   tokenManager,
		AuthService:    auth.NewService(cfg, logger, postgres.NewAuthRepository(pgClient), redisClient, tokenManager, newSmtp),
	}

	application.Run(ctx)
	application.Shutdown(ctx)
}
