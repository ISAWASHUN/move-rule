// @title          ゴミ分別クイズ API
// @version        1.0
// @description    自治体のゴミ分別をクイズ形式で学習するためのAPI
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes   http https
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/config"
	_ "github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/docs"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/infrastructure/middlewares"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/infrastructure/repository"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/infrastructure/repository/mysql"
	httpHandler "github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/interface/http"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/pkg"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/usecase"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.toml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(log)

	db, err := mysql.ConnectDB(cfg.MySQL, cfg.App.LogLevel)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	garbageItemRepo := repository.NewGarbageItemRepository(db)
	garbageCategoryRepo := repository.NewGarbageCategoryRepository(db)
	municipalityRepo := repository.NewMunicipalityRepository(db)

	quizUseCase := usecase.NewQuizUseCase(garbageItemRepo, garbageCategoryRepo, municipalityRepo)

	quizHandler := httpHandler.NewQuizHandler(quizUseCase)

	if cfg.App.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(pkg.Logger(log))
	r.Use(middlewares.CORS())

	api := r.Group("/api/v1")
	{
		quiz := api.Group("/quiz")
		{
			quiz.GET("/questions", quizHandler.GetQuestions)
			quiz.POST("/answer", quizHandler.PostAnswer)
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 404ハンドラー（JSON形式でエラーを返す）
	r.NoRoute(func(c *gin.Context) {
		pkg.HandleError(c, pkg.NewNotFoundError("エンドポイントが見つかりません", nil))
	})

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Graceful shutdownの設定
	go func() {
		log.Info("starting server", "address", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("server exited properly")
}
