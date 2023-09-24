package httpslog

import (
	"log/slog"
	"net/http"
	"time"
)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (rt RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt(r)
}

func RtWithDefaultLogger(rt http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return rt.RoundTrip(Default(r))
	})
}

func RoundTripper(inner http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		r = With(r,
			slog.Group("req",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.Int64("size", r.ContentLength),
			),
		)
		Debug(r, "Senf HTTP request")
		start := time.Now()
		resp, err := inner.RoundTrip(r)
		duration := time.Since(start)
		if err != nil {
			ErrorAttrs(resp.Request, "Outgoing HTTP request",
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
		LogAttrs(resp.Request, level, "End outgoing HTTP request",
			slog.Group("resp",
				slog.Int("status", resp.StatusCode),
				slog.Duration("duration", duration),
				slog.Int64("size", resp.ContentLength),
			),
		)
		return resp, nil
	})
}
