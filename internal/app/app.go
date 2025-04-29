package app

import (
	"auth-service/internal/config"
	"auth-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func NewApp() {
	config.ConnectDatabase()
	r := gin.Default()
	// r.Use(cors.Default())

	handler.RegisterRoutes(r)

	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello auth")
	})

	r.Run(":8080")
}
