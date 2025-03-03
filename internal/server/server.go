package server

// Pedro Henrique Goffi de Paulo
// REV (0, 1)

import (
	"bot-manager/internal/config"
	routes "bot-manager/internal/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Inicializa o servidor HTTP
func StartServer() {
	r := gin.Default()

	// Middleware para servir arquivos estáticos
	r.Static("/static", "./images")

	// Configurar rotas
	routes.SetupRoutes(r)
	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"success": "true"}) })
	r.POST("/", func(c *gin.Context) { c.JSON(200, gin.H{"success": "true"}) })

	// Rodar servidor
	port := config.GetServerPort()
	fmt.Println("Servidor rodando na porta", port)
	r.Run(":" + port)
}
