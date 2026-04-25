package adaptersinboundhttpmiddlewarelogger

import (
	"time"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
	"github.com/gin-gonic/gin"
)

const path = "adapters/inbound/http/middleware/logger"

func Logger(
	log portsoutboundlogging.Generic,
	on2xxMessage string,
	on4xxMessage string,
	on5xxMessage string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		const tag = path

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		meta := domainmodel.LogMeta{
			"status":     statusCode,
			"latency":    latency,
			"client_ip":  clientIP,
			"method":     method,
			"path":       path,
			"user_agent": c.Request.UserAgent(),
			"error":      errorMessage,
		}

		if statusCode >= 500 {
			log.Error(c, tag, on5xxMessage, meta)
		} else if statusCode >= 400 {
			log.Warn(c, tag, on4xxMessage, meta)
		} else {
			log.Info(c, tag, on2xxMessage, meta)
		}
	}
}
