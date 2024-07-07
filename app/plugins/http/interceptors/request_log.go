package interceptors

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"hr-system-go/app/plugins/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggingMiddleware logs the incoming HTTP request & duration.
func RequestLog(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		var payload interface{}
		if c.Request.Method == "GET" {
			payload = c.Request.URL.Query()
		} else {
			json.Unmarshal(requestBody, &payload)
		}

		c.Next()
		endTime := time.Now()

		logger.Info("Request",
			zap.String("Path", c.Request.URL.Path),
			zap.String("Method", c.Request.Method),
			zap.Any("Payload", payload),
			zap.Int("Status Code", c.Writer.Status()),
			zap.Any("Response", blw.body.String()),
			zap.Duration("Duration", endTime.Sub(startTime)),
		)
	}
}
