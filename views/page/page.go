package page

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/ancalabrese/gotth/views/components/head"
)

// ContentProviderFunc is a function that generates page-specific head metadata
// and content based on the incoming HTTP request.
// It returns the HeadViewModel, the main content component, and an optional error.
type ContentProviderFunc func(r *http.Request) (metadata head.HeadViewModel, content templ.Component, err error)

type WebPage struct {
	// Path is the path where this page is served (e.g. /about, /dashboard/products)
	Path string
	// Function to generate dynamic content and metadata
	ContentProvider ContentProviderFunc
}
