package routes

import (
	"bot-manager/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura as rotas da API
func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/tgimg", handlers.DownloadImageFromTelegram)
	}

}
