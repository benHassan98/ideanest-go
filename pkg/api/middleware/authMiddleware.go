package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func AuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		matched, err := regexp.MatchString("(^/organization)|(^/revoke-refresh-token)", c.Request.URL.Path)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !matched {
			c.Next()
			return
		}

		result, _ := redisClient.Get(ctx, strings.Split(c.Request.Header.Get("Authorization"), " ")[1]).Bool()

		if result {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}

}
