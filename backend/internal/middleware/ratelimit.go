package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	limiters    sync.Map
	cleanupOnce sync.Once
)

func getLimiter(ip string, r rate.Limit, burst int) *rate.Limiter {
	val, loaded := limiters.LoadOrStore(ip, &ipLimiter{
		limiter:  rate.NewLimiter(r, burst),
		lastSeen: time.Now(),
	})
	l := val.(*ipLimiter)
	if loaded {
		l.lastSeen = time.Now()
	}
	return l.limiter
}

// RateLimit ограничивает число запросов с одного IP.
// requests — сколько запросов разрешено за период window.
func RateLimit(requests int, window time.Duration) gin.HandlerFunc {
	r := rate.Every(window / time.Duration(requests))
	burst := requests

	cleanupOnce.Do(func() {
		go func() {
			for {
				time.Sleep(5 * time.Minute)
				limiters.Range(func(key, value any) bool {
					l := value.(*ipLimiter)
					if time.Since(l.lastSeen) > 10*time.Minute {
						limiters.Delete(key)
					}
					return true
				})
			}
		}()
	})

	return func(c *gin.Context) {
		if !getLimiter(c.ClientIP(), r, burst).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "слишком много запросов, подождите немного",
			})
			return
		}
		c.Next()
	}
}
