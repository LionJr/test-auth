package app

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"test-auth/internal/app/http/server"
	"test-auth/internal/config"
	"test-auth/internal/services/auth"
	"test-auth/pkg/token_manager"
)

type Application struct {
	Config         *config.AppConfig
	Logger         *zap.Logger
	PostgresClient *sqlx.DB
	RedisClient    *redis.Client
	AuthService    *auth.Service
	TokenManager   *token_manager.TokenManager
}

func (app *Application) Run(ctx context.Context) {
	httpServerErrCh := server.NewServer(
		ctx,
		app.Config,
		app.Logger,
		app.AuthService,
		app.TokenManager,
	)

	<-httpServerErrCh
}

func (app *Application) Shutdown(ctx context.Context) {
	app.Logger.Info("Shutdown database")
	_ = app.PostgresClient.Close()

	app.Logger.Info("Shutdown redis")
	_ = app.RedisClient.Shutdown(ctx)

	app.Logger.Info("Shutdown logger")
	_ = app.Logger.Sync()
}
