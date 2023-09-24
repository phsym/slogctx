package slogctx

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"
)

type logger struct {
	handler   slog.Handler
	addSource bool
	level     slog.Leveler
}

func get(ctx context.Context) *logger {
	if ctx == nil {
		return nil
	}
	l, _ := ctx.Value(_ctxKey{}).(*logger)
	return l
}

func (l *logger) save(ctx context.Context) context.Context {
	if l == nil || ctx == nil {
		return ctx
	}
	return &logCtx{ctx, l}
}

func newLogger(h slog.Handler) *logger {
	return &logger{h, false, nil}
}

func (l *logger) clone() *logger {
	if l == nil {
		return nil
	}
	clone := *l
	return &clone
}

func (l *logger) with(attrs ...slog.Attr) *logger {
	if l != nil && l.handler != nil {
		l = l.clone()
		l.handler = l.handler.WithAttrs(attrs)
	}
	return l
}

func (l *logger) withGroup(name string) *logger {
	if l != nil && l.handler != nil {
		l = l.clone()
		l.handler = l.handler.WithGroup(name)
	}
	return l
}

func (l *logger) withSource() *logger {
	if l != nil {
		l = l.clone()
		l.addSource = true
	}
	return l
}

func (l *logger) withLevel(level slog.Leveler) *logger {
	if l != nil {
		l = l.clone()
		l.level = level
	}
	return l
}

func (l *logger) shouldLog(ctx context.Context, level slog.Level) bool {
	return l != nil && l.handler != nil && (l.level != nil && level >= l.level.Level() || l.handler.Enabled(ctx, level))
}

func (l *logger) getCallerPc(skip uint) uintptr {
	if !AddSource && !l.addSource {
		return 0
	}
	pcs := [1]uintptr{0}
	runtime.Callers(2+int(skip), pcs[:])
	return pcs[0]
}

func (l *logger) log(ctx context.Context, callerSkip uint, level slog.Level, msg string, attrs ...any) {
	if l.shouldLog(ctx, level) {
		rec := slog.NewRecord(time.Now(), level, msg, l.getCallerPc(1+callerSkip))
		rec.Add(attrs...)
		l.handler.Handle(ctx, rec)
	}
}

func (l *logger) logAttrs(ctx context.Context, callerSkip uint, level slog.Level, msg string, attrs ...slog.Attr) {
	if l.shouldLog(ctx, level) {
		rec := slog.NewRecord(time.Now(), level, msg, l.getCallerPc(1+callerSkip))
		rec.AddAttrs(attrs...)
		l.handler.Handle(ctx, rec)
	}
}

func (l *logger) logf(ctx context.Context, callerSkip uint, level slog.Level, msg string, args ...any) {
	l.log(ctx, callerSkip+1, level, fmt.Sprintf(msg, args...))
}
