package ui

import (
	"context"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/roma-glushko/frens/internal/store"
)

// Server serves the embedded UI over HTTP.
type Server struct {
	addr   string
	server *http.Server
	logger *log.Logger
	store  store.Store
}

// NewServer creates a new UI server.
func NewServer(addr string, logger *log.Logger, s store.Store) *Server {
	return &Server{
		addr:   addr,
		logger: logger,
		store:  s,
	}
}

// Start starts the HTTP server and returns the actual address it's listening on.
func (s *Server) Start(ctx context.Context) (string, error) {
	assets, err := Assets()
	if err != nil {
		return "", fmt.Errorf("failed to load UI assets: %w", err)
	}

	mux := http.NewServeMux()

	// Register API routes
	if s.store != nil {
		api := NewAPI(s.store)
		api.RegisterRoutes(mux)
	}

	// Serve static files with SPA fallback
	fileServer := http.FileServer(http.FS(assets))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file directly
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		// Check if file exists
		f, err := assets.Open(path[1:]) // Remove leading slash
		if err != nil {
			// File not found, serve index.html for SPA routing
			if _, err := assets.Open("index.html"); err == nil {
				r.URL.Path = "/"
				fileServer.ServeHTTP(w, r)
				return
			}
			http.NotFound(w, r)
			return
		}
		f.Close()

		// Serve the actual file
		fileServer.ServeHTTP(w, r)
	})

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return "", fmt.Errorf("failed to listen on %s: %w", s.addr, err)
	}

	actualAddr := listener.Addr().String()

	s.server = &http.Server{
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := s.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Server error", "error", err)
		}
	}()

	return actualAddr, nil
}

// Stop gracefully shuts down the server.
func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.server.Shutdown(shutdownCtx)
}

// serveFile serves a file from the embedded filesystem.
func serveFile(w http.ResponseWriter, r *http.Request, assets fs.FS, path string) {
	f, err := assets.Open(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if seeker, ok := f.(fs.File); ok {
		if rs, ok := seeker.(interface {
			Read([]byte) (int, error)
			Seek(int64, int) (int64, error)
		}); ok {
			http.ServeContent(w, r, stat.Name(), stat.ModTime(), &readSeeker{rs})
			return
		}
	}

	http.ServeFile(w, r, path)
}

type readSeeker struct {
	rs interface {
		Read([]byte) (int, error)
		Seek(int64, int) (int64, error)
	}
}

func (r *readSeeker) Read(p []byte) (int, error) {
	return r.rs.Read(p)
}

func (r *readSeeker) Seek(offset int64, whence int) (int64, error) {
	return r.rs.Seek(offset, whence)
}
