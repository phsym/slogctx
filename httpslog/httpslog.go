package httpslog

import (
	"log/slog"
	"net/http"

	"github.com/phsym/slogctx"
)

var AddSource bool = false

func Disable(req *http.Request) *http.Request {
	return req.WithContext(slogctx.Disable(req.Context()))
}

func Default(req *http.Request) *http.Request {
	return req.WithContext(slogctx.Default(req.Context()))
}

func New(req *http.Request, h slog.Handler) *http.Request {
	return req.WithContext(slogctx.New(req.Context(), h))
}

func FromLogger(req *http.Request, l *slog.Logger) *http.Request {
	return req.WithContext(slogctx.FromLogger(req.Context(), l))
}

func Handler(req *http.Request) slog.Handler {
	return slogctx.Handler(req.Context())
}

func With(req *http.Request, attrs ...slog.Attr) *http.Request {
	return req.WithContext(slogctx.With(req.Context(), attrs...))
}

func WithGroup(req *http.Request, name string) *http.Request {
	return req.WithContext(slogctx.WithGroup(req.Context(), name))
}

func WithSource(req *http.Request) *http.Request {
	return req.WithContext(slogctx.WithSource(req.Context()))
}

func WithLevel(req *http.Request, level slog.Level) *http.Request {
	return req.WithContext(slogctx.WithLevel(req.Context(), level))
}

func Log(req *http.Request, level slog.Level, msg string, attrs ...any) {
	slogctx.Log(req.Context(), level, msg, attrs...)
}

func LogAttrs(req *http.Request, level slog.Level, msg string, attrs ...slog.Attr) {
	slogctx.LogAttrs(req.Context(), level, msg, attrs...)
}

func Logf(req *http.Request, level slog.Level, msg string, args ...any) {
	slogctx.Logf(req.Context(), level, msg, args...)
}

func Debug(req *http.Request, msg string, attrs ...any) {
	slogctx.Debug(req.Context(), msg, attrs...)
}

func Info(req *http.Request, msg string, attrs ...any) {
	slogctx.Info(req.Context(), msg, attrs...)
}

func Warn(req *http.Request, msg string, attrs ...any) {
	slogctx.Warn(req.Context(), msg, attrs...)
}

func Error(req *http.Request, msg string, attrs ...any) {
	slogctx.Error(req.Context(), msg, attrs...)
}

func DebugAttrs(req *http.Request, msg string, attrs ...slog.Attr) {
	slogctx.DebugAttrs(req.Context(), msg, attrs...)
}

func InfoAttrs(req *http.Request, msg string, attrs ...slog.Attr) {
	slogctx.InfoAttrs(req.Context(), msg, attrs...)
}

func WarnAttrs(req *http.Request, msg string, attrs ...slog.Attr) {
	slogctx.WarnAttrs(req.Context(), msg, attrs...)
}

func ErrorAttrs(req *http.Request, msg string, attrs ...slog.Attr) {
	slogctx.ErrorAttrs(req.Context(), msg, attrs...)
}

func Debugf(req *http.Request, msg string, attrs ...any) {
	slogctx.Debugf(req.Context(), msg, attrs...)
}

func Infof(req *http.Request, msg string, attrs ...any) {
	slogctx.Infof(req.Context(), msg, attrs...)
}

func Warnf(req *http.Request, msg string, attrs ...any) {
	slogctx.Warnf(req.Context(), msg, attrs...)
}

func Errorf(req *http.Request, msg string, attrs ...any) {
	slogctx.Errorf(req.Context(), msg, attrs...)
}
