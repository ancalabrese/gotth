package main

/*
import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ancalabrese/gotth/views/components/head"
)

// LoggingMiddleware is a simple example of a global middleware.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.RequestURI, time.Since(start))
	})
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Metadata for the Index Page
	indexHeadVM := head.NewHeadViewModel(
		head.WithPageCoreMetadata("Sample Home", "Welcome to our sample website built with Gotth!", "/"),
		head.WithKeywords([]string{"gotth", "sample", "homepage", "go", "templ"}),
		head.WithOpenGraph(
			"", "", // type, locale (use defaults)
			"/", "Sample Home OG Title", "OG description for sample home.",
			"https://placehold.co/1200x630/0779e4/ffffff?text=Sample+Home", "1200", "630", "Sample homepage OG image",
		),
	)

	pageDefinitions := []page.WebPage{
		{
			Path:     "/",
			Metadata: indexHeadVM,
			Content:  pages.IndexPageContent(),
		},
	}

	// Static assets for this sample app (e.g., global style.css)
	appStaticFS := server.NewStaticAssetFS(
		"/static",            // URL path where these assets will be served
		http.Dir("./static"), // Filesystem path relative to where main.go is run
	)

	// --- 4. Configure the WebServer ---
	// The `Layout` field expects a function matching `server.BaseLayoutFunc`
	// `views.BaseLayout` is a templ.Component that matches this signature.
	serverConfig := server.WebServerConfig{
		Layout:         views.BaseLayout, // Pass the BaseLayout component directly
		StaticAssetsFS: []server.StaticAssetFS{appStaticFS},
		GlobalMiddlewares: []func(http.Handler) http.Handler{
			LoggingMiddleware,
		},
	}

	// Create the underlying http.Server instance
	// You can configure timeouts, TLS, etc., here.
	httpServer := &http.Server{
		Addr:         ":8080", // Set the address here
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	srv, err := server.New(serverConfig, httpServer)
	if err != nil {
		log.Fatalf("Error creating web server: %v", err)
	}

	srv.RegisterPages(pageDefinitions)

	log.Println("Starting Gotth WebServer on http://localhost:8080")
	if err := srv.Start(ctx); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed.")
		} else {
			log.Fatalf("Server failed to start or closed unexpectedly: %v", err)
		}
	}
} */
