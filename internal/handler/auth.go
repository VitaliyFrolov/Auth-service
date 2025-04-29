package handler

import (
	"auth-service/internal/models"
	"auth-service/internal/service"
	"auth-service/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", func(c *gin.Context) {
			var req models.User
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if err := service.Register(req.Email, req.Password); err != nil {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "registered"})
		})

		auth.POST("/login", func(c *gin.Context) {
			var req models.User
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			accessToken, refreshToken, err := service.Login(req.Email, req.Password)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			})
		})

		auth.GET("/protected", func(c *gin.Context) {
			token := c.GetHeader("Authorization")
			if token == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
				return
			}

			userID, err := utils.ValidateToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"user_id": userID})
		})

		auth.POST("/refresh", func(c *gin.Context) {
			var req struct {
				RefreshToken string `json:"refresh_token"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			userID, err := utils.ValidateToken(req.RefreshToken)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
				return
			}

			accessToken, err := utils.GenerateAccessToken(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"access_token": accessToken,
			})
		})
	}
}
