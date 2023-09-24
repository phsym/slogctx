package httpslog

import (
	"log/slog"
	"net/http"
	"time"
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
			h.ServeHTTP(w, New(r, slog.Default().Handler()))
		})
	}
}

func WithLogger(logger slog.Handler) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, New(r, logger))
		})
	}
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = With(r,
			slog.Group("req",
				slog.String("method", r.Method),
				slog.String("uri", r.URL.String()),
				slog.String("ip", r.RemoteAddr),
				slog.String("agent", r.UserAgent()),
				slog.Int64("size", r.ContentLength),
			),
		)
		Debug(r, "Start incoming HTTP request")
		recorder := &statusRecorder{w, 0, 0}
		start := time.Now()
		h.ServeHTTP(recorder, r)
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
		LogAttrs(r, level, "End incoming HTTP request",
			slog.Group("resp",
				slog.Int("status", recorder.status),
				slog.Duration("duration", duration),
				slog.Uint64("size", recorder.size),
			),
		)
	})
}

// func WithLogger(logger slog.Handler) func(http.Handler) http.Handler {
// 	return func(h http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			h.ServeHTTP(w, r.WithContext(slogctx.New(r.Context(), logger)))
// 		})
// 	}
// }
