package middlewares

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	limiterMemory "github.com/ulule/limiter/v3/drivers/store/memory"
)

func ApiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		if apiKey != "my-secret-api-key" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			c.Abort()
			return
		}

		c.Next()
	}
}

var rateLimiter *limiter.Limiter

func init() {
	store := limiterMemory.NewStore()

	rate := limiter.Rate{
		Period: 1 + time.Minute,
		Limit:  5,
	}

	rateLimiter = limiter.New(store, rate)
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		context := context.Background()

		limiterCtx, err := rateLimiter.Get(context, ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter error"})
			c.Abort()
			return
		}

		if limiterCtx.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Writer.Header().Set("X-RateLimit-Limit", strconv.FormatInt(limiterCtx.Limit, 10))
		c.Writer.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(limiterCtx.Remaining, 10))
		c.Writer.Header().Set("X-RateLimit-Reset", strconv.FormatInt(limiterCtx.Reset, 10))

		c.Next()
	}
}

func CorsConfig() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length", "X-RateLimit-Limit", "X-RateLimit-Remaining"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

