// Copyright 2026 Roma Hlushko
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ui

import (
	"context"
	"errors"
	"fmt"
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

	api := NewAPI(s.store)
	api.RegisterRoutes(mux)

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

		if err := f.Close(); err != nil {
			s.logger.Warn("Failed to close file", "error", err)
		}

		// Serve the actual file
		fileServer.ServeHTTP(w, r)
	})

	lc := net.ListenConfig{}

	listener, err := lc.Listen(ctx, "tcp", s.addr)
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
		if err := s.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
