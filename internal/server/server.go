package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nihitdev/drop/internal/network"
	"github.com/nihitdev/drop/internal/share"
)

type StopReason int

const (
	StoppedByContext StopReason = iota
	StoppedAfterDownload
	StoppedAfterExpiry
)

type Config struct {
	Listener net.Listener
	Target   *share.Target
	Token    string
	Keep     bool
	Expiry   time.Duration
	OnClient func(string)
	OnSent   func(string)
}

type Server struct {
	cfg       Config
	http      *http.Server
	completed chan struct{}
	once      sync.Once
}

func New(cfg Config) *Server {
	s := &Server{cfg: cfg, completed: make(chan struct{})}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleDownload)
	s.http = &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	return s
}

func Run(ctx context.Context, cfg Config) (StopReason, error) {
	if cfg.Listener == nil || cfg.Target == nil || cfg.Token == "" {
		return StoppedByContext, fmt.Errorf("listener, target, and token are required")
	}
	s := New(cfg)
	serveErr := make(chan error, 1)
	go func() { serveErr <- s.http.Serve(cfg.Listener) }()

	var expiry <-chan time.Time
	var timer *time.Timer
	if cfg.Expiry > 0 {
		timer = time.NewTimer(cfg.Expiry)
		expiry = timer.C
		defer timer.Stop()
	}

	reason := StoppedByContext
	select {
	case <-ctx.Done():
	case <-s.completed:
		reason = StoppedAfterDownload
	case <-expiry:
		reason = StoppedAfterExpiry
	case err := <-serveErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return StoppedByContext, err
		}
		return StoppedByContext, nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := s.http.Shutdown(shutdownCtx); err != nil {
		_ = s.http.Close()
		return reason, fmt.Errorf("graceful shutdown: %w", err)
	}
	if err := <-serveErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return reason, err
	}
	return reason, nil
}

func (s *Server) Handler() http.Handler { return s.http.Handler }

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/d/"+s.cfg.Token {
		http.NotFound(w, r)
		return
	}

	file, err := os.Open(s.cfg.Target.Path)
	if err != nil {
		http.Error(w, "cannot open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "cannot inspect file", http.StatusInternalServerError)
		return
	}

	if s.cfg.OnClient != nil {
		s.cfg.OnClient(network.ClientIP(r))
	}
	w.Header().Set("Content-Disposition", contentDisposition(s.cfg.Target.DisplayName))
	recorder := &statusWriter{ResponseWriter: w, status: http.StatusOK}
	http.ServeContent(recorder, r, s.cfg.Target.DisplayName, stat.ModTime(), file)

	if r.Method == http.MethodGet && recorder.writeErr == nil && recorder.status >= 200 && recorder.status < 400 {
		if s.cfg.OnSent != nil {
			s.cfg.OnSent(s.cfg.Target.DisplayName)
		}
		if !s.cfg.Keep {
			s.once.Do(func() { close(s.completed) })
		}
	}
}

func contentDisposition(name string) string {
	name = strings.Map(func(r rune) rune {
		if r == '"' || r == '\\' || r == '\r' || r == '\n' || r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, name)
	if name == "" {
		name = "download"
	}
	return fmt.Sprintf(`attachment; filename="%s"`, name)
}

type statusWriter struct {
	http.ResponseWriter
	status   int
	writeErr error
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(p []byte) (int, error) {
	n, err := w.ResponseWriter.Write(p)
	if err != nil {
		w.writeErr = err
	}
	return n, err
}
