package slogctx

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"
)

func BenchmarkSlogBaseLogger(b *testing.B) {
	l := slog.New(slog.NewTextHandler(io.Discard, nil))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.WithGroup("group").Info("foobar")
	}
}

func BenchmarkSlogBase(b *testing.B) {
	h := slog.NewTextHandler(io.Discard, nil)
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := slog.NewRecord(time.Now(), slog.LevelInfo, "foobar", 0)
		h.WithGroup("group").Handle(ctx, r)
		// Info(ctx, "foobar")
	}
}

func BenchmarkSlogCtx(b *testing.B) {
	ctx := New(context.Background(), slog.NewTextHandler(io.Discard, nil))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info(WithGroup(ctx, "group"), "foobar")
	}
}
