package middleware

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"time"
)

func RateLimitMiddleware(max float64, ttl time.Duration) gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(max, &limiter.ExpirableOptions{DefaultExpirationTTL: ttl})

	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)

		if httpError != nil {
			c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}
