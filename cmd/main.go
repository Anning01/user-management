package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-management/internal/api"
	"user-management/internal/api/handlers"
	"user-management/internal/config"
	"user-management/internal/repository"
	"user-management/internal/service"
	"user-management/migrations"
	"user-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger.Init()

	// 连接数据库
	db, err := repository.NewDBConnection(&cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	// 运行数据库迁移
	if err := migrations.Migrate(db); err != nil {
		logger.Fatalf("Migration failed: %v", err)
	}

	// 初始化存储库
	userRepo := repository.NewUserRepository(db)
	articleRepo := repository.NewArticleRepository(db)

	// 初始化服务
	userService := service.NewUserService(userRepo)
	articleService := service.NewArticleService(articleRepo, userRepo)

	// 初始化处理器
	userHandler := handlers.NewUserHandler(userService, &cfg.JWT)
	articleHandler := handlers.NewArticleHandler(articleService)

	// 设置路由
	r := gin.Default()
	api.SetupRoutes(r, userHandler, articleHandler, &cfg.JWT)

	// 创建服务器
	srv := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        r,
		ReadTimeout:    cfg.Server.ReadTimeout * time.Second,
		WriteTimeout:   cfg.Server.WriteTimeout * time.Second,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	// 启动服务器（非阻塞）
	go func() {
		logger.Infof("Server is running on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", err)
	}

	logger.Info("Server exiting")
}
