package ginslog

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phsym/slogctx"
)

// func WithDefaultLogger() gin.HandlerFunc {
// 	return func(h http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			h.ServeHTTP(w, r.WithContext(slogctx.New(r.Context(), slog.Default().Handler())))
// 		})
// 	}
// }

// func WithLogger(logger slog.Handler)  gin.HandlerFunc {
// 	return func(h http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			h.ServeHTTP(w, r.WithContext(slogctx.New(r.Context(), logger)))
// 		})
// 	}
// }

func Middleware(h http.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		c := slogctx.With(r.Context(),
			slog.Group("req",
				slog.String("method", r.Method),
				slog.String("uri", r.URL.String()),
				slog.String("ip", r.RemoteAddr),
				slog.String("agent", r.UserAgent()),
				slog.Int64("size", r.ContentLength),
			),
		)
		slogctx.Info(c, "Start HTTP request")
		start := time.Now()
		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
		duration := time.Since(start)

		status := ctx.Writer.Status()
		level := slog.LevelInfo
		if status >= 400 {
			level = slog.LevelWarn
		} else if status >= 500 {
			level = slog.LevelError
		}
		slogctx.LogAttrs(ctx, level, "End HTTP request",
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
