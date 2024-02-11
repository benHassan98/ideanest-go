package pkg

import (
	"Ideanest/pkg/api/handlers"
	"Ideanest/pkg/api/middleware"
	"Ideanest/pkg/database/mongodb/repository"
	"Ideanest/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

func Init(config utils.Config) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := utils.ConnectToMongo(ctx, config.Database.Url)

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	redisClient := redis.NewClient(&redis.Options{Addr: config.Redis.Address})

	userRepository := repository.NewUserRepository(config, client)
	organizationRepository := repository.NewOrganizationRepository(config, client)

	userHandler := handlers.NewUserHandler(userRepository, redisClient)
	organizationHandler := handlers.NewOrganizationHandler(organizationRepository)

	router := gin.Default()

	router.Use(middleware.AuthMiddleware(redisClient))

	router.POST("/signup", userHandler.Signup)
	router.POST("/signin", userHandler.Signin)
	router.POST("/refresh-token", userHandler.RefreshToken)
	router.POST("/revoke-refresh-token", userHandler.RevokeRefreshToken)

	router.POST("/organization", organizationHandler.Create)
	router.GET("/organization/:organization_id", organizationHandler.FindById)
	router.GET("/organization", organizationHandler.FindAll)
	router.PUT("/organization/:organization_id", organizationHandler.Update)
	router.POST("/organization/:organization_id/invite", organizationHandler.InviteMember)
	router.DELETE("/organization/:organization_id", organizationHandler.DeleteById)

	if err := router.Run(":" + config.Server.Port); err != nil {
		return err
	}
	return nil
}
