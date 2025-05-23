package page

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/ancalabrese/gotth/internal/views/components/head"
)

type WebPage struct {
	// Path is the path where this page is served (e.g. /about, /dashboard/products)
	Path string
	// SEO and metadata for this specific page
	Metadata head.HeadViewModel
	// The templ.Component for the main body of this page
	Content templ.Component
	// HandlerFunc use when content needs more complex logic before rendering
	HandlerFunc http.HandlerFunc
}
