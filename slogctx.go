package slogctx

import (
	"context"
	"fmt"
	"log/slog"
)

var AddSource bool = false

func Disable(ctx context.Context) context.Context {
	return New(ctx, nil)
}

func Default(ctx context.Context) context.Context {
	return New(ctx, slog.Default().Handler())
}

func New(ctx context.Context, h slog.Handler) context.Context {
	if ctx == nil {
		return nil
	}
	return newLogger(h).save(ctx)
}

func FromLogger(ctx context.Context, l *slog.Logger) context.Context {
	if ctx == nil {
		return nil
	}
	return newLogger(l.Handler()).save(ctx)
}

func Handler(ctx context.Context) slog.Handler {
	return get(ctx).handler
}

func With(ctx context.Context, attrs ...slog.Attr) context.Context {
	return get(ctx).with(attrs...).save(ctx)
}

func WithGroup(ctx context.Context, name string) context.Context {
	return get(ctx).withGroup(name).save(ctx)
}

func WithSource(ctx context.Context) context.Context {
	return get(ctx).withSource().save(ctx)
}

func WithLevel(ctx context.Context, level slog.Level) context.Context {
	return get(ctx).withLevel(level).save(ctx)
}

func Log(ctx context.Context, level slog.Level, msg string, attrs ...any) {
	get(ctx).log(ctx, 1, level, msg, attrs...)
}

func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	get(ctx).logAttrs(ctx, 1, level, msg, attrs...)
}

func Logf(ctx context.Context, level slog.Level, msg string, args ...any) {
	get(ctx).logf(ctx, 1, level, fmt.Sprintf(msg, args...))
}

func Debug(ctx context.Context, msg string, attrs ...any) {
	get(ctx).log(ctx, 1, slog.LevelDebug, msg, attrs...)
}

func Info(ctx context.Context, msg string, attrs ...any) {
	get(ctx).log(ctx, 1, slog.LevelInfo, msg, attrs...)
}

func Warn(ctx context.Context, msg string, attrs ...any) {
	get(ctx).log(ctx, 1, slog.LevelWarn, msg, attrs...)
}

func Error(ctx context.Context, msg string, attrs ...any) {
	get(ctx).log(ctx, 1, slog.LevelError, msg, attrs...)
}

func DebugAttrs(ctx context.Context, msg string, attrs ...slog.Attr) {
	get(ctx).logAttrs(ctx, 1, slog.LevelDebug, msg, attrs...)
}

func InfoAttrs(ctx context.Context, msg string, attrs ...slog.Attr) {
	get(ctx).logAttrs(ctx, 1, slog.LevelInfo, msg, attrs...)
}

func WarnAttrs(ctx context.Context, msg string, attrs ...slog.Attr) {
	get(ctx).logAttrs(ctx, 1, slog.LevelWarn, msg, attrs...)
}

func ErrorAttrs(ctx context.Context, msg string, attrs ...slog.Attr) {
	get(ctx).logAttrs(ctx, 1, slog.LevelError, msg, attrs...)
}

func Debugf(ctx context.Context, msg string, attrs ...any) {
	get(ctx).logf(ctx, 1, slog.LevelDebug, msg, attrs...)
}

func Infof(ctx context.Context, msg string, attrs ...any) {
	get(ctx).logf(ctx, 1, slog.LevelInfo, msg, attrs...)
}

func Warnf(ctx context.Context, msg string, attrs ...any) {
	get(ctx).logf(ctx, 1, slog.LevelWarn, msg, attrs...)
}

func Errorf(ctx context.Context, msg string, attrs ...any) {
	get(ctx).logf(ctx, 1, slog.LevelError, msg, attrs...)
}
