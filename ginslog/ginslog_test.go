package ginslog_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/phsym/slogctx/ginslog"
)

func TestGinsLog(t *testing.T) {
	engine := gin.New()
	buf := bytes.Buffer{}
	logger := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	engine.Use(ginslog.WithLogger(logger), ginslog.AccessLog())

	engine.GET("/test", func(ctx *gin.Context) {
		ginslog.Warn(ctx, "FOOBAR")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	engine.ServeHTTP(httptest.NewRecorder(), req)

	println(buf.String())

	dec := json.NewDecoder(&buf)
	entry := LogEntry{}
	if err := dec.Decode(&entry); err != nil {
		t.Fatalf("failed to decode log: %s", err.Error())
	}
	entry.AssertField(t, "msg", "Start incoming HTTP request")
	entry.AssertField(t, "level", "DEBUG")
	reqEntry := entry.AssertObject(t, "req")
	reqEntry.AssertField(t, "method", "GET")
	reqEntry.AssertField(t, "uri", "/test")

	clear(entry)
	if err := dec.Decode(&entry); err != nil {
		t.Fatalf("failed to decode log: %s", err.Error())
	}
	entry.AssertField(t, "msg", "FOOBAR")
	entry.AssertField(t, "level", "WARN")
	reqEntry = entry.AssertObject(t, "req")
	reqEntry.AssertField(t, "method", "GET")
	reqEntry.AssertField(t, "uri", "/test")

	clear(entry)
	if err := dec.Decode(&entry); err != nil {
		t.Fatalf("failed to decode log: %s", err.Error())
	}
	entry.AssertField(t, "msg", "End incoming HTTP request")
	entry.AssertField(t, "level", "INFO")
	reqEntry = entry.AssertObject(t, "req")
	reqEntry.AssertField(t, "method", "GET")
	reqEntry.AssertField(t, "uri", "/test")
	respEntry := entry.AssertObject(t, "resp")
	respEntry.AssertField(t, "status", float64(200))
	// t.Fail()
}

type LogEntry (map[string]any)

func (log LogEntry) AssertFieldExists(t *testing.T, key string) any {
	if f, ok := log[key]; ok {
		return f
	}
	t.Fatalf("Missing field %q", key)
	return nil
}

func (log LogEntry) AssertField(t *testing.T, key string, expect any) {
	f := log.AssertFieldExists(t, key)
	if f != expect {
		t.Fatalf("Unexpectd field %q. Expected %+v, got %+v", key, expect, f)
	}
}

func (log LogEntry) AssertObject(t *testing.T, key string) LogEntry {
	f := log.AssertFieldExists(t, key)
	obj, ok := f.(map[string]any)
	if !ok {
		t.Fatalf("Expected an object, got %T", f)
		return nil
	}
	return obj
}
