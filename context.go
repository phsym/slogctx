package slogctx

import (
	"context"
	"time"
)

type _ctxKey struct{}

type logCtx struct {
	parent context.Context
	logger *logger
}

// Deadline implements context.Context.
func (c *logCtx) Deadline() (deadline time.Time, ok bool) {
	return c.parent.Deadline()
}

// Done implements context.Context.
func (c *logCtx) Done() <-chan struct{} {
	return c.parent.Done()
}

// Err implements context.Context.
func (c *logCtx) Err() error {
	return c.parent.Err()
}

// Value implements context.Context.
func (c *logCtx) Value(key any) any {
	if _, ok := key.(_ctxKey); ok {
		return c.logger
	}
	return c.parent.Value(key)
}

var _ context.Context = (*logCtx)(nil)
