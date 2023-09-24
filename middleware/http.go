package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/phsym/slogctx"
)

type statusRecorder struct {
	rw     http.ResponseWriter
	status int
	size   uint64
}

// Header implements http.ResponseWriter.
func (s *statusRecorder) Header() http.Header {
	return s.rw.Header()
}

// Write implements http.ResponseWriter.
func (s *statusRecorder) Write(b []byte) (int, error) {
	if s.status <= 0 {
		s.status = http.StatusOK
		s.rw.WriteHeader(s.status)
	}
	n, err := s.rw.Write(b)
	s.size += uint64(n)
	return n, err
}

// WriteHeader implements http.ResponseWriter.
func (s *statusRecorder) WriteHeader(statusCode int) {
	if s.status > 0 {
		return
	}
	s.status = statusCode
	s.rw.WriteHeader(statusCode)
}

var _ http.ResponseWriter = (*statusRecorder)(nil)

func WithDefaultLogger() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r.WithContext(slogctx.New(r.Context(), slog.Default().Handler())))
		})
	}
}

func WithLogger(logger slog.Handler) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r.WithContext(slogctx.New(r.Context(), logger)))
		})
	}
}

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := slogctx.With(r.Context(),
			slog.Group("req",
				slog.String("method", r.Method),
				slog.String("uri", r.URL.String()),
				slog.String("ip", r.RemoteAddr),
				slog.String("agent", r.UserAgent()),
				slog.Int64("size", r.ContentLength),
			),
		)
		slogctx.Info(ctx, "Start HTTP request")
		recorder := &statusRecorder{w, 0, 0}
		start := time.Now()
		h.ServeHTTP(recorder, r.WithContext(ctx))
		duration := time.Since(start)
		level := slog.LevelInfo
		if recorder.status >= 400 {
			level = slog.LevelWarn
		} else if recorder.status >= 500 {
			level = slog.LevelError
		}
		if recorder.status <= 0 {
			recorder.status = http.StatusOK
		}
		slogctx.LogAttrs(ctx, level, "End HTTP request",
			slog.Group("resp",
				slog.Int("status", recorder.status),
				slog.Duration("duration", duration),
				slog.Uint64("size", recorder.size),
			),
		)
	})
}

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (rt RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt(r)
}

func RtWithDefaultLogger(rt http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return rt.RoundTrip(r.WithContext(slogctx.Default(r.Context())))
	})
}

// func WithLogger(logger slog.Handler) func(http.Handler) http.Handler {
// 	return func(h http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			h.ServeHTTP(w, r.WithContext(slogctx.New(r.Context(), logger)))
// 		})
// 	}
// }

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
