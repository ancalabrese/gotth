package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/ancalabrese/gotth"
	"github.com/ancalabrese/gotth/example/middleware"
	"github.com/ancalabrese/gotth/example/views"
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

	// Static assets for this sample app (e.g., global style.css)
	appStaticFS := gotth.NewStaticAssetFS(
		"/static",                  // URL path where these assets will be served
		http.Dir("./static/dist/"), // Filesystem path relative to where main.go is run
	)

	cfg := gotth.WebServerConfig{
		StaticAssetsFS: []gotth.StaticAssetFS{appStaticFS},
		GlobalMiddlewares: []func(http.Handler) http.Handler{
			middleware.GottherName,
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

	webServer, err := gotth.New(cfg, httpServer)
	if err != nil {
		panic(err)
	}

	webServer.ServeContent("/", func(r *http.Request) (metadata head.HeadViewModel, content templ.Component, err error) {
		indexHeadVM := head.NewHeadViewModel(
			head.WithHTMX(""),
			head.WithPageCoreMetadata(
				"Gotth: Get your site online fast | Go + Templ + Tailwind + HTMX",
				"Get your Go websites online fast with Gotth! Leverages Templ, Tailwind CSS, and HTMX for rapid, SEO-friendly development of modern web applications.",
				"/index"),
			head.WithKeywords([]string{"Gotth", "Go web server", "Templ", "Tailwind CSS", "HTMX", "Go templ", "Fast Go websites", "Rapid web development Go", "Go SEO", "Go web starter kit", "Go Tailwind HTMX", "Full-stack Go"}),
			head.WithFavicon("/static/gotth.svg", "image/png"),
			head.WithKeywords([]string{"gotth", "sample", "homepage", "go", "templ"}),
			head.WithOpenGraph(
				"", "", // type, locale (use defaults)
				"/", "Sample Home OG Title", "OG description for sample home.",
				"https://placehold.co/1200x630/0779e4/ffffff?text=Sample+Home", "1200", "630", "Sample homepage OG image",
			),
			head.WithStylesheet("/static/style.css", "", "", ""),
		)

		name := middleware.GetGottherName(r.Context())
		if name == "" {
			content = views.Home()
		} else {
			content = views.HomeWithName(name)
		}
		return indexHeadVM, content, nil
	})

	webServer.Start(ctx)
}
