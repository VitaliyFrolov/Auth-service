package service

import (
	"auth-service/internal/config"
	"auth-service/internal/models"
	"auth-service/internal/utils"
	"fmt"
)

func Register(username, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := models.User{
		Username: username,
		Password: hashedPassword,
	}

	return config.DB.Create(&user).Error
}

func Login(username, password string) (string, string, error) {
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", "", err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", fmt.Errorf("invalid credentials")
	}
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
