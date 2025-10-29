package api

import (
	"github.com/Anning01/user-management/internal/api/handlers"
	"github.com/Anning01/user-management/internal/api/middleware"
	"github.com/Anning01/user-management/internal/config"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	userHandler *handlers.UserHandler,
	articleHandler *handlers.ArticleHandler,
	jwtConfig *config.JWTConfig,
) {
	// 公开路由
	public := r.Group("/api/v1")
	{
		// 用户相关
		public.POST("/users/register", userHandler.Register)
		public.POST("/users/login", userHandler.Login)

		// 文章相关（公开访问的）
		public.GET("/articles", articleHandler.ListArticles)
		public.GET("/articles/:id", articleHandler.GetArticle)
	}

	// 需要认证的路由
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(jwtConfig))
	{
		// 用户相关
		protected.GET("/users/me", userHandler.GetCurrentUser)
		protected.PUT("/users/me", userHandler.UpdateCurrentUser)
		protected.DELETE("/users/me", userHandler.DeleteCurrentUser)

		// 文章相关
		protected.POST("/articles", articleHandler.CreateArticle)
		protected.PUT("/articles/:id", articleHandler.UpdateArticle)
		protected.DELETE("/articles/:id", articleHandler.DeleteArticle)
		protected.GET("/users/me/articles", articleHandler.ListMyArticles)
	}
}
