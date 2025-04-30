package service

import (
	"auth-service/internal/config"
	"auth-service/internal/models"
	"auth-service/internal/utils"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Register(email, password string) (uint, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}

func Login(email, password string) (string, string, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
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
