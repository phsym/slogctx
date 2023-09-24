package httpslog

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/phsym/slogctx"
)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (rt RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt(r)
}

func RtWithDefaultLogger(rt http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return rt.RoundTrip(r.WithContext(slogctx.Default(r.Context())))
	})
}

func RoundTripper(inner http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		ctx := slogctx.With(r.Context(),
			slog.Group("req",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.Int64("size", r.ContentLength),
			),
		)
		// slogctx.Info(ctx, "Start HTTP request")
		start := time.Now()
		resp, err := inner.RoundTrip(r.WithContext(ctx))
		duration := time.Since(start)
		if err != nil {
			slogctx.ErrorAttrs(ctx, "Outgoing HTTP request",
				slog.Group("resp",
					slog.Any("err", err),
				),
			)
			return nil, err
		}
		level := slog.LevelInfo
		if resp.StatusCode >= 400 {
			level = slog.LevelWarn
		} else if resp.StatusCode >= 500 {
			level = slog.LevelError
		}
		slogctx.LogAttrs(ctx, level, "End HTTP request",
			slog.Group("resp",
				slog.Int("status", resp.StatusCode),
				slog.Duration("duration", duration),
				slog.Int64("size", resp.ContentLength),
			),
		)
		return resp, nil
	})
}
