![logo](/assets/Gotth.svg)
# Gotth: Fast Web App Development with Go, Templ, Tailwind & HTMX 🚀

**Gotth** offers a structured way to serve dynamic UIs with `templ`, handle all that important HTML metadata, and integrate smoothly with modern frontend tools like Tailwind CSS and HTMX.

"Ship fast, fail fast" - As we all know.

## Core Features

* **`gotth.WebServer`: Your Web Server Foundation**:
    * A ready-to-go HTTP server. You can easily plug in global middlewares and tell it where your static assets (CSS, JS, images) are.
    * Comes with graceful shutdown, because smooth deployments are happy deployments.
    * Uses the standard `http.ServeMux` for routing when you use `ServeContent` directly.

* **Define your content with `ContentProviderFunc`**:
    * This is how you define what each page shows. A `ContentProviderFunc` returns your page-specific metadata (for the `<head>`) and the main `templ.Component` for the body.


* **Static File Serving (`StaticAssetFS`)**:
    * Straightforward configuration for serving your static files (CSS, JavaScript, images).

* **`head.HeadViewModel`: Comprehensive HTML `<head>` Management**:
    * Managing the HTML `<head>` can be a pain, right? `HeadViewModel` makes it much easier.
    * It uses a clean functional options pattern (those `With...` functions) to set things up.
    * **Smart Defaults**: It can automatically figure out Open Graph and Twitter card metadata from your basic page info if you don't specify them.
    * **Supports a Ton of Stuff**:
        * Core metadata: title, description, canonical URL, keywords, author.
        * Icons: Favicons, Apple touch icons.
        * Microsoft stuff: MS tile and browser configuration.
        * Social Sharing: Open Graph tags (type, locale, URL, title, description, image, etc.) and Twitter Card tags.
        * Theming: `theme-color`, `apple-mobile-web-app-status-bar-style`, `color-scheme`.
        * Custom Fonts: `FontLink` for your web fonts.
        * Stylesheets: `StylesheetLink` with media, integrity, and crossorigin attributes.
        * Header Scripts: `ScriptLink` with async, defer, type, integrity, and crossorigin.
        * **Easy Frontend Libs**: Quickly include HTMX, HTMX Preload Extension, and Alpine.js from CDNs (SRI included, and you can use your own URLs) or create your own with Header Scripts.
        * Analytics: Just drop in your Google Analytics tag ID.
        * Custom Meta Tags: If you need anything else, you can add it.

* **`head.JSONLDNode`: Structured Data with JSON-LD**:
    * Want to give search engines more detailed info? `JSONLDNode` helps you define JSON-LD objects.
    * It handles the specific JSON-LD keywords (`@context`, `@id`, `@type`) and custom properties correctly when marshaling to JSON.
    * You can add this to your `HeadViewModel` using `WithPreparedJSONLD` (if you have a pre-made JSON string) or `WithJSONLD` (to marshal a `JSONLDNode` object directly – *note: the `WithJSONLD` in `head/viewmodel.go` is a placeholder and needs full marshaling logic*).

* **`viewmodel.NewViewModel`: Generic ViewModel Builder**:
    * A handy generic function `NewViewModel[T any](data ...ViewModelData[T]) *T`. It uses functional options to initialize pretty much any view model struct you define. Less boilerplate is good!

* **Session Management (`middlewares` package)**:
    * If you need user sessions, Gotth provides the basics.
    * `SessionStore` interface: Abstract away your session storage (e.g., database, Redis). You implement `ExchangeSessionIDForUser` and `InvalidateSession`.
    * `SessionCheck` middleware: Checks for a session cookie, validates it with your `SessionStore`, gets the user, and puts the user info into the request context. Can handle required or optional sessions.
    * `InvalidateSession` middleware: Clears the session from your `SessionStore` and removes the cookie.
    * `GetUser(ctx context.Context)`: A helper to easily get the user from the request context.

* **Example App**:
    * Check out `example/cmd/main.go`. It's a working example showing how to use Gotth with `templ`, Tailwind, and HTMX. You'll see:
        * Server setup.
        * Static file serving.
        * How to write a `ContentProviderFunc`.
        * Using `HeadViewModel`.
        * Adding custom middleware.

## Getting Started

1.  **Install**:
    ```bash
    go get github.com/ancalabrese/gotth
    ```

2.  **Basic Usage**:

    ```go
    package main

    import (
    	"context"
    	"io"
    	"log"
    	"net/http"

    	"[github.com/a-h/templ](https://github.com/a-h/templ)"
    	"[github.com/ancalabrese/gotth](https://github.com/ancalabrese/gotth)"
    	"[github.com/ancalabrese/gotth/views/components/head](https://github.com/ancalabrese/gotth/views/components/head)"
    	// You'll likely have your own layout package, e.g.:
    	// import "your-project/views/layout"
    )

    // Your page's content provider function
    func myPageContentProvider(r *http.Request) (head.HeadViewModel, templ.Component, error) {
    	headVM := head.NewHeadViewModel(
    		head.WithPageCoreMetadata("My Page Title", "A cool description of my page.", r.URL.Path),
    		head.WithStylesheet("/static/css/style.css", "", "", ""),
    		head.WithHTMX(""), // Use default HTMX CDN, easy peasy
    	)

    	// Your page's main content (a templ.Component)
    	pageContent := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
    		// This would be a more complex .templ component in a real app
    		_, err := io.WriteString(w, "<div><h1>Welcome!</h1><p>Page content lives here.</p></div>")
    		return err
    	})

    	return headVM, pageContent, nil
    }

    // A simplified layout component example
    // In a real app, this would be a proper templ component rendering <html>, <head>, <body>, etc.
    func basicLayout(headVM head.HeadViewModel, content templ.Component) templ.Component {
        return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
            log.Println("Rendering basic layout with title:", headVM.Core.Title) // Example of using headVM
            return content.Render(ctx, w)
        })
    }

    func main() {
    	ctx, cancel := context.WithCancel(context.Background())
    	defer cancel()

    	// Setup static asset serving
    	staticFS := gotth.NewStaticAssetFS("/static", http.Dir("./path/to/your/static/files")) // Update this path!

    	cfg := gotth.WebServerConfig{
    		StaticAssetsFS: []gotth.StaticAssetFS{staticFS},
    		// Add any global middlewares here:
    		// GlobalMiddlewares: []func(http.Handler) http.Handler{ /* your middleware */ },
    	}

    	httpServer := &http.Server{
    		Addr: ":8080",
    	}

    	webServer, err := gotth.New(cfg, httpServer)
    	if err != nil {
    		log.Fatalf("Failed to create web server: %v", err)
    	}

    	// Register your page handler
    	webServer.ServeContent("/my-page", myPageContentProvider)

    	log.Println("Server starting on http://localhost:8080 ✨")
    	if err := webServer.Start(ctx); err != nil {
    		log.Fatalf("Server failed to start: %v", err)
    	}
    }
    ```

## Key Concepts

* **`ContentProviderFunc`**: This is central to how Gotth organizes page generation. It keeps your page-specific data and component logic separate from routing.
* **`head.HeadViewModel`**: A structured way to build your HTML `<head>`. Good for SEO, social sharing, and managing your assets.
* **Functional Options**: You'll see this pattern a lot (e.g., `With...` functions). It makes configuring components cleaner and more explicit.
* **Middleware**: Use global server middleware or the provided session middleware to build up your request processing.
* **`templ` Integration**: Gotth is built to work hand-in-hand with `templ` for server-side rendering your HTML.

## Future Directions & Router Compatibility

Gotth currently uses `http.ServeMux` internally. However, the core page rendering logic (`ContentProviderFunc`, `HeadViewModel`) is designed to be adaptable. You could make it work with other routers like `chi` or `gorilla/mux` by:

1.  Creating a function that takes your `ContentProviderFunc` (and maybe a layout function) and returns a standard `http.HandlerFunc`.
2.  Registering these handlers with your router of choice.
3.  Exposing the static file serving as an `http.Handler` too.

This would make Gotth more flexible for different Go web project setups.

## Contributing
Contributions are welcome!
Please open an issue to discuss bigger changes or submit a pull request for fixes and small improvements.

