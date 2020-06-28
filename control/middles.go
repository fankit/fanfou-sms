package control

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
	"tsmsrv/utils"
)

func AuthMiddleWare() gin.HandlerFunc {

	return func(ctx *gin.Context) {

	}
}

//日志中间件
func LoggerMidderWare() gin.HandlerFunc {
	utils.Logger.Sync(3 * time.Second)

	return func(ctx *gin.Context) {
		nt := time.Now()

		ctx.Next()
		latency := time.Since(nt)
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()
		path := ctx.Request.URL.Path

		switch {
		case statusCode >= 400 && statusCode <= 499:
			utils.Logger.Log.Error("http",
				zap.String("method", method),
				zap.Int("code", statusCode),
				zap.String("clientIP", clientIP),
				zap.String("latency", latency.String()),
				zap.String("path", path),
				zap.String("desc", ctx.Errors.String()),//ctx.Errors.String()),
			)
		case statusCode >= 500:
			utils.Logger.Log.Warn("http",
				zap.String("method", method),
				zap.Int("code", statusCode),
				zap.String("clientIP", clientIP),
				zap.String("latency", latency.String()),
				zap.String("path", path),
				zap.String("desc", ctx.Errors.String()),
			)
		default:
			utils.Logger.Log.Info("http",
				zap.String("method", method),
				zap.Int("code", statusCode),
				zap.String("clientIP", clientIP),
				zap.String("latency", latency.String()),
				zap.String("path", path),
				zap.String("desc", ctx.Errors.String()),
			)
		}

	}
}
func SourceNotFound(c *gin.Context)  {
	c.String(http.StatusNotFound, "Source Not Found")
}