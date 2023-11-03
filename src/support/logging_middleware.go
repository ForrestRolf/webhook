package support

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()
		duration := endTime.Sub(startTime)

		log.SetFormatter(&log.JSONFormatter{})
		log.WithFields(log.Fields{
			"client_ip": ctx.ClientIP(),
			"duration":  duration,
			"method":    ctx.Request.Method,
			"status":    ctx.Writer.Status(),
			"referrer":  ctx.Request.Referer(),
		}).Info(ctx.Request.RequestURI)

		ctx.Next()
	}
}
