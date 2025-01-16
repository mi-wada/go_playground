package testslog_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"testing"
)

var mu sync.Mutex

func captureSlog(t *testing.T, f func()) []byte {
	t.Helper()

	mu.Lock()
	defer mu.Unlock()
	original := slog.Default()

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	slog.SetDefault(logger)
	defer slog.SetDefault(original)

	f()

	return buf.Bytes()
}

func myFunc() {
	slog.Info("myFunc called", slog.String("key", "value"))
}

func TestMyFunc(t *testing.T) {
	t.Parallel()

	for i := range 10 {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			logBytes := captureSlog(t, func() {
				myFunc()
			})

			var m map[string]interface{}
			if err := json.Unmarshal(logBytes, &m); err != nil {
				t.Fatalf("json.Unmarshal() = %v; want nil", err)
			}

			if got, want := m["msg"], "myFunc called"; got != want {
				t.Errorf("msg = %v; want %v", got, want)
			}
			if got, want := m["level"], slog.LevelInfo.String(); got != want {
				t.Errorf("level = %v; want %v", got, want)
			}
			if got, want := m["key"], "value"; got != want {
				t.Errorf("key = %v; want %v", got, want)
			}
		})
	}
}
