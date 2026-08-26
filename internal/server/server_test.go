package server

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nihitdev/drop/internal/share"
)

func testTarget(t *testing.T, contents string) *share.Target {
	t.Helper()
	path := filepath.Join(t.TempDir(), "example.txt")
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	return &share.Target{Path: path, DisplayName: "example.txt", Size: int64(len(contents))}
}

func TestHandlerRequiresExactTokenAndServesFile(t *testing.T) {
	s := New(Config{Target: testTarget(t, "hello"), Token: "secret", Keep: true})

	bad := httptest.NewRecorder()
	s.Handler().ServeHTTP(bad, httptest.NewRequest(http.MethodGet, "/d/wrong", nil))
	if bad.Code != http.StatusNotFound {
		t.Fatalf("invalid token status = %d, want 404", bad.Code)
	}

	good := httptest.NewRecorder()
	s.Handler().ServeHTTP(good, httptest.NewRequest(http.MethodGet, "/d/secret", nil))
	if good.Code != http.StatusOK || good.Body.String() != "hello" {
		t.Fatalf("valid response = (%d, %q), want (200, %q)", good.Code, good.Body.String(), "hello")
	}
}

func TestOneShotStopsAfterDownload(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	result := make(chan StopReason, 1)
	errs := make(chan error, 1)
	target := testTarget(t, "hello")
	go func() {
		reason, err := Run(context.Background(), Config{Listener: listener, Target: target, Token: "secret"})
		result <- reason
		errs <- err
	}()

	resp, err := http.Get("http://" + listener.Addr().String() + "/d/secret")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	select {
	case reason := <-result:
		if err := <-errs; err != nil {
			t.Fatal(err)
		}
		if reason != StoppedAfterDownload {
			t.Fatalf("stop reason = %v, want download", reason)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server did not stop after download")
	}
}

func TestKeepStillExpires(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	started := time.Now()
	reason, err := Run(context.Background(), Config{
		Listener: listener, Target: testTarget(t, "hello"), Token: "secret", Keep: true, Expiry: 30 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	if reason != StoppedAfterExpiry {
		t.Fatalf("stop reason = %v, want expiry", reason)
	}
	if time.Since(started) < 20*time.Millisecond {
		t.Fatal("server expired too early")
	}
}
