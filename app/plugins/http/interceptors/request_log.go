package interceptors

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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
	return func(ctx *gin.Context) {
		startTime := time.Now()

		var requestBody []byte
		if ctx.Request.Body != nil {
			requestBody, _ = io.ReadAll(ctx.Request.Body)
		}

		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		var payload interface{}
		if ctx.Request.Method == http.MethodGet {
			payload = ctx.Request.URL.Query()
		} else {
			_ = json.Unmarshal(requestBody, &payload)
		}

		ctx.Next()
		endTime := time.Now()

		logger.Info("Request",
			zap.String("Path", ctx.Request.URL.Path),
			zap.String("Method", ctx.Request.Method),
			zap.Any("Payload", payload),
			zap.Int("Status Code", ctx.Writer.Status()),
			zap.Any("Response", blw.body.String()),
			zap.Duration("Duration", endTime.Sub(startTime)),
		)
	}
}
