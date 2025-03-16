package ratelimiter

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"time"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379", // Replace with your Redis server address
})

const (
	REQUEST_LIMIT = 100 // Max requests per second
	TIME_WINDOW   = 1   // Time window in seconds
)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		username := c.Request.Header.Get("username")
		key := fmt.Sprintf("ratelimit:%s", username)

		// Increment request count
		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			c.Abort()
			return
		}

		// Set expiration time for the key if it's a new counter
		if count == 1 {
			redisClient.Expire(ctx, key, time.Duration(TIME_WINDOW)*time.Second)
		}

		// If the user exceeds the limit, block the request
		if count > REQUEST_LIMIT {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests, slow down!"})
			c.Abort()
			return
		}

		// Proceed with the request
		c.Next()
	}
}
