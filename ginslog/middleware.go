package ginslog

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func WithDefaultLogger() gin.HandlerFunc {
	return WithLogger(slog.Default().Handler())
}

func WithLogger(logger slog.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx = New(ctx, logger)
		ctx.Next()
	}
}

func AccessLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx = With(ctx,
			slog.Group("req",
				slog.String("method", ctx.Request.Method),
				slog.String("uri", ctx.Request.URL.String()),
				slog.String("ip", ctx.ClientIP()),
				slog.String("agent", ctx.Request.UserAgent()),
				slog.Int64("size", ctx.Request.ContentLength),
				slog.String("path", ctx.FullPath()),
			),
		)
		Debug(ctx, "Start incoming HTTP request")
		start := time.Now()
		ctx.Next()
		duration := time.Since(start)

		status := ctx.Writer.Status()
		level := slog.LevelInfo
		if status >= 400 {
			level = slog.LevelWarn
		} else if status >= 500 {
			level = slog.LevelError
		}
		LogAttrs(ctx, level, "End incoming HTTP request",
			slog.Group("resp",
				slog.Int("status", status),
				slog.Duration("duration", duration),
				slog.Int("size", ctx.Writer.Size()),
			),
		)
	}
}
