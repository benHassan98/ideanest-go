package handlers

import (
	"Ideanest/pkg/database/mongodb/models"
	"Ideanest/pkg/database/mongodb/repository"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type UserHandler struct {
	repository  repository.UserRepository
	redisClient *redis.Client
}

type refreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result)
}

func NewUserHandler(repository repository.UserRepository, redisClient *redis.Client) UserHandler {

	return UserHandler{repository: repository, redisClient: redisClient}

}

func (u UserHandler) Signup(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, err := u.repository.InsertOne(ctx, user)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user created"})

}

func (u UserHandler) Signin(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		return
	}

	err := u.repository.FindOne(ctx, user.Email, user.Password)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken := generateRandomString(10), generateRandomString(10)+generateRandomString(5)

	u.redisClient.Set(ctx, accessToken, true, 0)
	u.redisClient.Set(ctx, refreshToken, true, 0)

	c.JSON(http.StatusOK, gin.H{"message": "signin successful", "access_token": accessToken, "refresh_token": refreshToken})

}

func (u UserHandler) RefreshToken(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var token refreshToken

	if err := c.BindJSON(&token); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result, err := u.redisClient.Get(ctx, token.RefreshToken).Bool()

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if result {

		u.redisClient.Del(ctx, token.RefreshToken)
		u.redisClient.Del(ctx, strings.Split(c.Request.Header.Get("Authorization"), " ")[1])

		accessToken, refreshToken := generateRandomString(10), generateRandomString(10)+generateRandomString(5)

		u.redisClient.Set(ctx, accessToken, true, 0)
		u.redisClient.Set(ctx, refreshToken, true, 0)

		c.JSON(http.StatusOK, gin.H{"message": "issued new refresh token successfully", "access_token": accessToken, "refresh_token": refreshToken})

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid refresh token"})
	}

}

func (u UserHandler) RevokeRefreshToken(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var token refreshToken

	if err := c.BindJSON(&token); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	u.redisClient.Del(ctx, strings.Split(c.Request.Header.Get("Authorization"), " ")[1])
	u.redisClient.Del(ctx, token.RefreshToken)

	c.JSON(http.StatusOK, gin.H{"message": "signout successful"})

}
