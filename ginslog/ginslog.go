package ginslog

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/phsym/slogctx/httpslog"
)

var AddSource bool = false

func Disable(ctx *gin.Context) *gin.Context {
	ctx = ctx.Copy()
	ctx.Request = (httpslog.Disable(ctx.Request))
	return ctx
}

func Default(ctx *gin.Context) *gin.Context {
	ctx = ctx.Copy()
	ctx.Request = (httpslog.Default(ctx.Request))
	return ctx
}

func New(ctx *gin.Context, h slog.Handler) *gin.Context {
	ctx = ctx.Copy()
	ctx.Request = (httpslog.New(ctx.Request, h))
	return ctx
}

func FromLogger(ctx *gin.Context, l *slog.Logger) *gin.Context {
	ctx = ctx.Copy()
	ctx.Request = (httpslog.FromLogger(ctx.Request, l))
	return ctx
}

func Handler(ctx *gin.Context) slog.Handler {
	return httpslog.Handler(ctx.Request)
}

func With(ctx *gin.Context, attrs ...slog.Attr) *gin.Context {
	ctx = ctx.Copy()
	ctx.Request = (httpslog.With(ctx.Request, attrs...))
	return ctx
}

func WithGroup(ctx *gin.Context, name string) *gin.Context {
	ctx = ctx.Copy()
	ctx.Request = (httpslog.WithGroup(ctx.Request, name))
	return ctx
}

func WithSource(ctx *gin.Context) *gin.Context {
	ctx = ctx.Copy()
	ctx.Request = (httpslog.WithSource(ctx.Request))
	return ctx
}

func WithLevel(ctx *gin.Context, level slog.Level) *gin.Context {
	ctx = ctx.Copy()
	ctx.Request = (httpslog.WithLevel(ctx.Request, level))
	return ctx
}

func Log(ctx *gin.Context, level slog.Level, msg string, attrs ...any) {
	httpslog.Log(ctx.Request, level, msg, attrs...)
}

func LogAttrs(ctx *gin.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	httpslog.LogAttrs(ctx.Request, level, msg, attrs...)
}

func Logf(ctx *gin.Context, level slog.Level, msg string, args ...any) {
	httpslog.Logf(ctx.Request, level, msg, args...)
}

func Debug(ctx *gin.Context, msg string, attrs ...any) {
	httpslog.Debug(ctx.Request, msg, attrs...)
}

func Info(ctx *gin.Context, msg string, attrs ...any) {
	httpslog.Info(ctx.Request, msg, attrs...)
}

func Warn(ctx *gin.Context, msg string, attrs ...any) {
	httpslog.Warn(ctx.Request, msg, attrs...)
}

func Error(ctx *gin.Context, msg string, attrs ...any) {
	httpslog.Error(ctx.Request, msg, attrs...)
}

func DebugAttrs(ctx *gin.Context, msg string, attrs ...slog.Attr) {
	httpslog.DebugAttrs(ctx.Request, msg, attrs...)
}

func InfoAttrs(ctx *gin.Context, msg string, attrs ...slog.Attr) {
	httpslog.InfoAttrs(ctx.Request, msg, attrs...)
}

func WarnAttrs(ctx *gin.Context, msg string, attrs ...slog.Attr) {
	httpslog.WarnAttrs(ctx.Request, msg, attrs...)
}

func ErrorAttrs(ctx *gin.Context, msg string, attrs ...slog.Attr) {
	httpslog.ErrorAttrs(ctx.Request, msg, attrs...)
}

func Debugf(ctx *gin.Context, msg string, attrs ...any) {
	httpslog.Debugf(ctx.Request, msg, attrs...)
}

func Infof(ctx *gin.Context, msg string, attrs ...any) {
	httpslog.Infof(ctx.Request, msg, attrs...)
}

func Warnf(ctx *gin.Context, msg string, attrs ...any) {
	httpslog.Warnf(ctx.Request, msg, attrs...)
}

func Errorf(ctx *gin.Context, msg string, attrs ...any) {
	httpslog.Errorf(ctx.Request, msg, attrs...)
}
