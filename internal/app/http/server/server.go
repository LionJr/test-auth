package server

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	errch "github.com/proxeter/errors-channel"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"test-auth/internal/config"
	"test-auth/internal/services/auth"
	"test-auth/pkg/token_manager"
)

type Server struct {
	Logger       *zap.Logger
	Config       *config.AppConfig
	AuthService  *auth.Service
	tokenManager *token_manager.TokenManager
}

func NewServer(
	ctx context.Context,
	config *config.AppConfig,
	logger *zap.Logger,
	authService *auth.Service,
	tokenManager *token_manager.TokenManager,
) <-chan error {
	return errch.Register(func() error {
		return (&Server{
			Logger:       logger,
			Config:       config,
			AuthService:  authService,
			tokenManager: tokenManager,
		}).Start(ctx)
	})
}

func (s *Server) Start(ctx context.Context) error {
	h := s.initHandlers()

	server := http.Server{
		Handler:        h,
		Addr:           ":" + strconv.Itoa(s.Config.HTTP.Port),
		MaxHeaderBytes: s.Config.HTTP.MaxHeaderBytes,
		ReadTimeout:    s.Config.HTTP.ReadTimeout,
		WriteTimeout:   s.Config.HTTP.WriteTimeout,
	}

	s.Logger.Info(
		"Server running",
		zap.String("host", s.Config.HTTP.Host),
		zap.Int("port", s.Config.HTTP.Port),
	)

	select {
	case err := <-errch.Register(server.ListenAndServe):
		s.Logger.Info("Shutdown server", zap.String("by", "error"), zap.Error(err))
		return server.Shutdown(ctx)
	case <-ctx.Done():
		s.Logger.Info("Shutdown server", zap.String("by", "context.Done"))
		return server.Shutdown(ctx)
	}
}

func (s *Server) initHandlers() *gin.Engine {
	// Init gin handler
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET, POST, OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:8080"},
		MaxAge:           86400,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	api := router.Group("/api")
	authRouter := api.Group("/auth")
	authRouter.POST("/sign-in", s.AuthService.SignInHandler)
	authRouter.POST("/check-code", s.AuthService.CheckCodeHandler)
	authRouter.POST("/refresh", s.AuthService.RefreshHandler)

	return router
}
