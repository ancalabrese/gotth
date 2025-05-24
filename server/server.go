package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/a-h/templ"
	"github.com/ancalabrese/gotth/views/components/head"
	"github.com/ancalabrese/gotth/views/page"
)

type StaticAssetFS struct {
	// URL path for the static assets (e.g., "/static").
	urlPath string
	assetFS http.FileSystem
}

func NewStaticAssetFS(url string, fs http.FileSystem) StaticAssetFS {
	return StaticAssetFS{
		urlPath: url,
		assetFS: fs,
	}
}

// BaseLayoutFunc is the signature for the function/component that wraps page content.
type BaseLayoutFunc func(headVM head.HeadViewModel, pageContent templ.Component) templ.Component

// WebServerConfig holds the config for WebServer
type WebServerConfig struct {
	// Optional: FSs for global static assets (CSS/JS/Assets etc)
	StaticAssetsFS []StaticAssetFS
	// The BaseLayout component func
	Layout BaseLayoutFunc
	// Middlewares globally applied
	GlobalMiddlewares []func(http.Handler) http.Handler
}

// WebServer handles HTTP requests and serves configured web pages
type WebServer struct {
	config     WebServerConfig
	httpServer *http.Server
	mux        *http.ServeMux // Using standard library ServeMux for simplicity
}

// New creates a new WebServer.
func New(cfg WebServerConfig, s *http.Server) (*WebServer, error) {
	if s == nil {
		s = defaultServer()
	}

	mux := http.NewServeMux()
	// Setup global static file serving if configured
	for _, fsConfig := range cfg.StaticAssetsFS {
		if fsConfig.assetFS != nil && fsConfig.urlPath != "" {
			urlPath := fsConfig.urlPath
			if !strings.HasPrefix(urlPath, "/") {
				urlPath = "/" + urlPath
			}

			// Path for Handle needs a trailing slash to match subpaths.
			// Prefix for StripPrefix should NOT have the trailing slash if the Handle path does.
			servePath := urlPath
			if !strings.HasSuffix(servePath, "/") {
				servePath += "/"
			}

			stripPath := strings.TrimSuffix(urlPath, "/")

			mux.Handle(servePath, http.StripPrefix(stripPath, http.FileServer(fsConfig.assetFS)))
			fmt.Printf("Serving static assets from URL path '%s'\n", servePath)
		}
	}

	return &WebServer{
		httpServer: s,
		config:     cfg,
		mux:        mux,
	}, nil
}

// RegisterPage adds a page to be served.
func (ws *WebServer) RegisterPage(p page.WebPage) {
	if p.Path == "" || p.ContentProvider == nil {
		fmt.Printf("Skipping registration of page with empty path or no ContentProvider\n")
		return
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headVM, pageContent, err := p.ContentProvider(r)
		if err != nil {
			// TODO: Handle the error appropriately (e.g., log it, show a generic error page)
			// allow the ContentProviderFunc to also suggest an HTTP status code
			fmt.Fprintf(os.Stderr, "Error in ContentProvider for %s: %v\n", p.Path, err)
			return
		}

		// Create the full page component by wrapping the page's content with the base layout
		fullPageComponent := ws.config.Layout(headVM, pageContent)

		// Set content type and render
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = fullPageComponent.Render(r.Context(), w) // Pass request context
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error rendering page %s: %v\n", p.Path, err)
			// On rendering error return HTTP error. Any other error should be an error message
			// in the rendered page. TODO: better error handling
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	fmt.Printf("Registering page at path: %s\n", p.Path)
	ws.mux.Handle(p.Path, handler)
}

// RegisterPages registers multiple pages.
func (ws *WebServer) RegisterPages(pages []page.WebPage) {
	for _, p := range pages {
		ws.RegisterPage(p)
	}
}

// Start initializes and runs the HTTP server.
// Cancelling the context will stop the server
func (ws *WebServer) Start(ctx context.Context) error {
	var finalHandler http.Handler = ws.mux
	// Apply in reverse
	for i := len(ws.config.GlobalMiddlewares) - 1; i >= 0; i-- {
		finalHandler = ws.config.GlobalMiddlewares[i](finalHandler)
	}
	ws.httpServer.Handler = finalHandler

	fmt.Printf("WebServer starting on %s\n", ws.httpServer.Addr)
	ws.httpServer.Handler = ws.mux

	errChan := make(chan error, 1)
	go func() {
		if err := ws.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("ListenAndServe failed: %w", err)
		}
		close(errChan)
	}()

	// Wait for an error or a shutdown signal
	select {
	case err := <-errChan:
		return err
	case <-ws.gracefulShutdownContext(ctx).Done():
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := ws.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}
		fmt.Println("WebServer gracefully stopped")
		return nil
	}
}

func (ws *WebServer) gracefulShutdownContext(ctx context.Context) context.Context {
	cancellableCtx, cancel := context.WithCancel(ctx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer signal.Stop(sigChan)
		<-sigChan
		cancel()
	}()
	return cancellableCtx
}

func defaultServer() *http.Server {
	return &http.Server{
		Addr:              ":8080",
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}
}
