package ginslog

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func WithDefaultLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx = New(ctx, slog.Default().Handler())
		ctx.Next()
	}
}

// func WithLogger(logger slog.Handler)  gin.HandlerFunc {
// 	return func(h http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			h.ServeHTTP(w, r.WithContext(slogctx.New(r.Context(), logger)))
// 		})
// 	}
// }

func Middleware() gin.HandlerFunc {
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
				slog.Uint64("size", uint64(ctx.Writer.Size())),
			),
		)
	}
}

// func WithLogger(logger slog.Handler) func(http.Handler) http.Handler {
// 	return func(h http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			h.ServeHTTP(w, r.WithContext(slogctx.New(r.Context(), logger)))
// 		})
// 	}
// }
