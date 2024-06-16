package logs

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

func Init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func buildAttr(c *gin.Context) []slog.Attr {
	if c == nil {
		return nil
	}
	return []slog.Attr{
		slog.String("method", c.Request.Method),
		slog.String("path", c.Request.URL.Path),
		slog.String("ip", c.ClientIP()),
	}
}

func Info(c *gin.Context, msg string) {
	slog.Info(msg, buildAttr(c))
}

func Error(c *gin.Context, msg string) {
	slog.Error(msg, buildAttr(c))
}

func Warn(c *gin.Context, msg string) {
	slog.Warn(msg, buildAttr(c))
}

func Debug(c *gin.Context, msg string) {
	slog.Debug(msg, buildAttr(c))
}
