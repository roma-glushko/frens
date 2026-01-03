// Copyright 2025 Roma Hlushko
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

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/log"
	jctx "github.com/roma-glushko/frens/internal/context"
	"github.com/roma-glushko/frens/internal/ui"
	"github.com/urfave/cli/v2"
)

var ServeCommand = &cli.Command{
	Name:  "serve",
	Usage: "Start the web UI server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "addr",
			Aliases: []string{"a"},
			Value:   "127.0.0.1:8080",
			Usage:   "Address to listen on (host:port)",
		},
		&cli.BoolFlag{
			Name:    "open",
			Aliases: []string{"o"},
			Value:   false,
			Usage:   "Open the UI in the default browser",
		},
	},
	Action: func(c *cli.Context) error {
		addr := c.String("addr")
		openBrowser := c.Bool("open")

		logger := log.NewWithOptions(os.Stderr, log.Options{
			Level: log.InfoLevel,
		})

		appCtx := jctx.FromCtx(c.Context)
		if appCtx == nil {
			return fmt.Errorf("failed to get app context")
		}

		server := ui.NewServer(addr, logger, appCtx.Store)

		actualAddr, err := server.Start(c.Context)
		if err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}

		url := fmt.Sprintf("http://%s", actualAddr)
		logger.Info("Frens UI is running", "url", url)

		if openBrowser {
			if err := openURL(url); err != nil {
				logger.Warn("Failed to open browser", "error", err)
			}
		}

		// Wait for interrupt signal
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.Info("Shutting down server...")

		if err := server.Stop(c.Context); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}

		logger.Info("Server stopped")
		return nil
	},
}

// openURL opens the specified URL in the default browser.
func openURL(url string) error {
	var cmd string
	var args []string

	switch {
	case isWSL():
		cmd = "cmd.exe"
		args = []string{"/c", "start", url}
	case isMacOS():
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}

	return runCommand(cmd, args...)
}

func isWSL() bool {
	_, err := os.Stat("/proc/version")
	if err != nil {
		return false
	}
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return contains(string(data), "microsoft") || contains(string(data), "WSL")
}

func isMacOS() bool {
	return os.Getenv("OSTYPE") == "darwin" || fileExists("/usr/bin/open")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func runCommand(name string, args ...string) error {
	cmd := newCommand(name, args...)
	return cmd.Start()
}

// newCommand creates a new exec.Cmd - this is a simple wrapper for testing
func newCommand(name string, args ...string) *command {
	return &command{name: name, args: args}
}

type command struct {
	name string
	args []string
}

func (c *command) Start() error {
	proc, err := os.StartProcess(c.name, append([]string{c.name}, c.args...), &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})
	if err != nil {
		return err
	}
	return proc.Release()
}
