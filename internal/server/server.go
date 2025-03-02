package server

import (
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

	// Rodar servidor
	port := "8080"
	fmt.Println("Servidor rodando na porta", port)
	r.Run(":" + port)
}
