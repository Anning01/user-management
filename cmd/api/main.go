package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Anning01/user-management/internal/api"
	"github.com/Anning01/user-management/internal/api/handlers"
	"github.com/Anning01/user-management/internal/config"
	"github.com/Anning01/user-management/internal/repository"
	"github.com/Anning01/user-management/internal/service"
	"github.com/Anning01/user-management/migrations"
	"github.com/Anning01/user-management/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载 .env 文件（如果存在）
	// 如果文件不存在，不会报错，会继续使用系统环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	} else {
		log.Println("Loaded .env file successfully")
	}

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 打印配置信息（用于调试，不显示敏感信息）
	log.Printf("Configuration loaded:")
	log.Printf("  Server Port: %s", cfg.Server.Port)
	log.Printf("  Database Host: %s", cfg.Database.Host)
	log.Printf("  Database Port: %s", cfg.Database.Port)
	log.Printf("  Database Name: %s", cfg.Database.Name)
	log.Printf("  Database Username: %s", cfg.Database.Username)
	if cfg.Database.Password == "" {
		log.Printf("  Database Password: (empty) - connecting without password")
	} else {
		log.Printf("  Database Password: (set, length: %d)", len(cfg.Database.Password))
	}
	if cfg.JWT.SecretKey == "" {
		log.Printf("  ⚠️  WARNING: JWT Secret Key is empty! This is insecure for production!")
	} else {
		log.Printf("  JWT Secret Key: (set, length: %d)", len(cfg.JWT.SecretKey))
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
