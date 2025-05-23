// Package head provides view models and components for HTML head section.
package head

// HeadViewModel is the primary model for the Head templ component.
// Instantiate via NewHeadViewModel and functional options.
type HeadViewModel struct {
	// General Site/Application Information
	Name string // Application name, used for application-name meta, OG site_name, etc.

	// Core Page Metadata (SEO and page identity)
	Metadata PageMetadata

	// Favicons and Touch Icons
	FaviconPath        string // Path to the main favicon (e.g., /static/favicon.ico or /static/image.png)
	FaviconType        string // MIME type of the favicon (e.g., image/x-icon, image/png)
	AppleTouchIconPath string // Path to the Apple touch icon (e.g., /static/apple-touch-icon.png)

	// Microsoft
	MsTileColor         string // Tile color for Windows pinned sites (e.g., #FFFFFF)
	MsBrowserConfigPath string // Path to browserconfig.xml (if used)
	MsStartURL          string // Start URL for MS applications (defaults to "/")

	// Open Graph (Social Sharing - Facebook, LinkedIn, etc.) -
	// OgURL, OgTitle, OgDescription, OgImage, etc., are populated from PageMetadata
	OgType   string // OpenGraph type (defaults to "website")
	OgLocale string // OpenGraph locale (defaults to "en_US")

	// Twitter Card (Social Sharing - Twitter)
	// TwitterURL, TwitterTitle, TwitterDescription, TwitterImage are populated from PageMetadata
	TwitterCardType      string // Twitter card type (defaults to "summary_large_image")
	TwitterSiteHandle    string // Twitter handle of the website (e.g., "@YourSite")
	TwitterCreatorHandle string // Twitter handle of the content creator (e.g., "@AuthorHandle")

	// Theming and PWA-like Behavior
	ThemeColor          string // Theme color for browser UI (Android)
	AppleStatusBarColor string // Status bar style for iOS web apps (e.g., "black-translucent")
	ColorScheme         string // Supported color schemes (e.g., "light dark")

	// Common JavaScript Libraries (Booleans for inclusion, Paths for location)
	IncludeHTMX        bool
	HTMXPath           string // Defaults to a common CDN or local path
	IncludeHTMXPreload bool
	HTMXPreloadPath    string // Defaults to a common CDN or local path
	IncludeAlpineJS    bool
	AlpineJSPath       string // Defaults to a common CDN or local path

	// Custom Assets
	Fonts         []FontLink       // List of fonts to link
	Stylesheets   []StylesheetLink // List of CSS stylesheets
	HeaderScripts []ScriptLink     // List of JavaScript files for the head

	// Analytics
	IsAnalyticsEnabled bool
	MeasuramentID      string // e.g., Google Analytics Measurement ID

	// Structured Data
	PreparedJSONLD string // Pre-marshaled JSON-LD string

	// Miscellaneous
	CustomMetaTags map[string]string // For any other arbitrary meta tags
}

// PageMetadata contains detailed metadata for a specific page.
type PageMetadata struct {
	Title       string // Page title
	Description string // Page description (for meta description and social)
	URL         string // Canonical URL of the page

	// Optional - recommended for better SEO and user experience
	Author         string   // Author of the content
	Keywords       []string // SEO Keywords
	SchemaImageURL string   // Main image URL for basic schema.org itemprop="image"

	ViewPort string

	// Open Graph Specifics (fallbacks from Title, Description, URL, SchemaImageURL if not explicitly set)
	OgURL         string
	OgTitle       string
	OgDescription string
	OgImage       string // URL of the OpenGraph image
	OgImageWidth  string // Width of the OG image (e.g., "1200")
	OgImageHeight string // Height of the OG image (e.g., "630")
	OgImageAlt    string // Alt text for OG image

	// Twitter Specifics (fallbacks from OG tags or main metadata if not explicitly set)
	TwitterTitle       string
	TwitterDescription string
	TwitterImage       string // URL of the Twitter card image
	TwitterImageAlt    string // Alt text for Twitter image
}

// FontLink defines a font to be loaded.
type FontLink struct {
	Href        string // Full URL to the font CSS or font file
	CrossOrigin bool   // Add crossorigin attribute if true
}

// StylesheetLink defines a CSS stylesheet.
type StylesheetLink struct {
	Href        string
	Media       string // Optional (e.g., "screen", "print")
	Integrity   string // Optional, for Subresource Integrity
	CrossOrigin string // Optional ("anonymous", "use-credentials")
}

// ScriptLink defines a JavaScript file for the head.
type ScriptLink struct {
	Src         string
	IsAsync     bool
	IsDefer     bool
	Type        string // Optional (e.g., "module")
	Integrity   string // Optional
	CrossOrigin string // Optional ("anonymous", "use-credentials")
}

// Option defines a function that sets a field in HeadViewModel.
type Option func(*HeadViewModel)

// NewHeadViewModel creates a new HeadViewModel with default values when possible.
// Use Options to set values.
// For a meaningful <head>, callers should provide options for:
//   - WithPageMetadata (setting Title, Description, URL)
//   - WithName (if application-name is desired)
//
// Other fields have defaults or are optional.
func NewHeadViewModel(opts ...Option) HeadViewModel {
	vm := HeadViewModel{
		Metadata: PageMetadata{
			ViewPort: "width=device-width, initial-scale=1.0",
		},
		MsStartURL:      "/",
		OgType:          "website",
		OgLocale:        "en_US",
		TwitterCardType: "summary_large_image",
		// Default JS library paths (can be overridden by options)
		HTMXPath:        "/static/js/htmx.min.js",
		HTMXPreloadPath: "https://unpkg.com/htmx-ext-preload@2.0.1/preload.js",
		AlpineJSPath:    "https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js",
		CustomMetaTags:  make(map[string]string),
	}

	for _, opt := range opts {
		opt(&vm)
	}

	// --- Post-processing: Consolidate fallbacks after all options are applied ---
	// This ensures that if specific OG/Twitter values weren't set by dedicated options,
	// they fall back to general metadata or other social metadata.

	// Open Graph fallbacks
	if vm.Metadata.OgURL == "" {
		vm.Metadata.OgURL = vm.Metadata.URL
	}
	if vm.Metadata.OgTitle == "" {
		vm.Metadata.OgTitle = vm.Metadata.Title
	}
	if vm.Metadata.OgDescription == "" {
		vm.Metadata.OgDescription = vm.Metadata.Description
	}
	if vm.Metadata.TwitterImage == "" {
		vm.Metadata.TwitterImage = vm.Metadata.OgImage
	}
	if vm.Metadata.TwitterImageAlt == "" && vm.Metadata.OgImageAlt != "" {
		vm.Metadata.TwitterImageAlt = vm.Metadata.OgImageAlt // Fallback Twitter alt to OG alt
	} else if vm.Metadata.TwitterImageAlt == "" && vm.Metadata.TwitterTitle != "" { // Basic alt from title if nothing else
		vm.Metadata.TwitterImageAlt = "Image for " + vm.Metadata.TwitterTitle
	}

	return vm
}

// WithName sets the application name.
func WithName(name string) Option {
	return func(vm *HeadViewModel) { vm.Name = name }
}

// WithPageCoreMetadata sets the most essential page metadata fields.
// Title, Description, and CanonicalURL should always be provided.
func WithPageCoreMetadata(title, description, canonicalURL string) Option {
	return func(vm *HeadViewModel) {
		vm.Metadata.Title = title
		vm.Metadata.Description = description
		vm.Metadata.URL = canonicalURL
	}
}

// WithAuthor sets the page author.
func WithAuthor(author string) Option {
	return func(vm *HeadViewModel) { vm.Metadata.Author = author }
}

// WithKeywords sets SEO keywords.
func WithKeywords(keywords []string) Option {
	return func(vm *HeadViewModel) { vm.Metadata.Keywords = keywords }
}

// WithViewport overrides the default viewport setting.
func WithViewport(viewport string) Option {
	return func(vm *HeadViewModel) {
		if viewport != "" {
			vm.Metadata.ViewPort = viewport
		}
	}
}

// WithSchemaImageURL sets the main image for basic schema.org itemprop="image".
func WithSchemaImageURL(url string) Option {
	return func(vm *HeadViewModel) { vm.Metadata.SchemaImageURL = url }
}

// WithFavicon sets the path and type for the website's favicon.
func WithFavicon(path, favType string) Option {
	return func(vm *HeadViewModel) {
		vm.FaviconPath = path
		vm.FaviconType = favType
	}
}

// WithAppleTouchIcon sets the path for the Apple touch icon.
func WithAppleTouchIcon(path string) Option {
	return func(vm *HeadViewModel) { vm.AppleTouchIconPath = path }
}

// WithMicrosoftOptions configures Microsoft-specific meta tags.
func WithMicrosoftOptions(tileColor, browserConfigPath, startURL string) Option {
	return func(vm *HeadViewModel) {
		if tileColor != "" {
			vm.MsTileColor = tileColor
		}
		if browserConfigPath != "" {
			vm.MsBrowserConfigPath = browserConfigPath
		}
		if startURL != "" {
			vm.MsStartURL = startURL
		} // Override default "/" if specified
	}
}

// WithOpenGraph configures Open Graph tags. Empty strings for individual params will not override existing values
// unless explicitly handled (here, they will override if not empty).
func WithOpenGraph(ogType, ogLocale, ogURL, ogTitle, ogDescription, ogImage, imgWidth, imgHeight, imgAlt string) Option {
	return func(vm *HeadViewModel) {
		if ogType != "" {
			vm.OgType = ogType
		}
		if ogLocale != "" {
			vm.OgLocale = ogLocale
		}
		if ogURL != "" {
			vm.Metadata.OgURL = ogURL
		}
		if ogTitle != "" {
			vm.Metadata.OgTitle = ogTitle
		}
		if ogDescription != "" {
			vm.Metadata.OgDescription = ogDescription
		}
		if ogImage != "" {
			vm.Metadata.OgImage = ogImage
		}
		if imgWidth != "" {
			vm.Metadata.OgImageWidth = imgWidth
		}
		if imgHeight != "" {
			vm.Metadata.OgImageHeight = imgHeight
		}
		if imgAlt != "" {
			vm.Metadata.OgImageAlt = imgAlt
		}
	}
}

// WithTwitterCard configures Twitter Card tags.
func WithTwitterCard(cardType, siteHandle, creatorHandle, title, description, image, imageAlt string) Option {
	return func(vm *HeadViewModel) {
		if cardType != "" {
			vm.TwitterCardType = cardType
		}
		if siteHandle != "" {
			vm.TwitterSiteHandle = siteHandle
		}
		if creatorHandle != "" {
			vm.TwitterCreatorHandle = creatorHandle
		}
		if title != "" {
			vm.Metadata.TwitterTitle = title
		}
		if description != "" {
			vm.Metadata.TwitterDescription = description
		}
		if image != "" {
			vm.Metadata.TwitterImage = image
		}
		if imageAlt != "" {
			vm.Metadata.TwitterImageAlt = imageAlt
		}
	}
}

// WithThemeing configures theme-color, Apple status bar, and color-scheme.
func WithThemeing(themeColor, appleStatusColor, colorScheme string) Option {
	return func(vm *HeadViewModel) {
		if themeColor != "" {
			vm.ThemeColor = themeColor
		}
		if appleStatusColor != "" {
			vm.AppleStatusBarColor = appleStatusColor
		}
		if colorScheme != "" {
			vm.ColorScheme = colorScheme
		}
	}
}

// WithAnalytics enables or disables analytics and sets the measurement ID.
func WithAnalytics(enabled bool, measurementID string) Option {
	return func(vm *HeadViewModel) {
		vm.IsAnalyticsEnabled = enabled
		vm.MeasuramentID = measurementID
	}
}

// WithPreparedJSONLD sets the pre-marshaled JSON-LD string.
func WithPreparedJSONLD(jsonLD string) Option {
	return func(vm *HeadViewModel) { vm.PreparedJSONLD = jsonLD }
}

// WithJSONLD marshals the jsonLD object and sets the JSON_LD value.
func WithJSONLD(jsonLD JsonLD) Option {
	return func(hvm *HeadViewModel) {}
}

// WithFont adds a font link to the list of fonts.
func WithFont(href string, crossOrigin bool) Option {
	return func(vm *HeadViewModel) {
		vm.Fonts = append(vm.Fonts, FontLink{Href: href, CrossOrigin: crossOrigin})
	}
}

// WithStylesheet adds a CSS stylesheet link.
func WithStylesheet(href, media, integrity, crossOrigin string) Option {
	return func(vm *HeadViewModel) {
		vm.Stylesheets = append(vm.Stylesheets, StylesheetLink{
			Href: href, Media: media, Integrity: integrity, CrossOrigin: crossOrigin,
		})
	}
}

// WithHeaderScript adds a JavaScript link to be included in the head.
func WithHeaderScript(src string, isAsync, isDefer bool, scriptType, integrity, crossOrigin string) Option {
	return func(vm *HeadViewModel) {
		vm.HeaderScripts = append(vm.HeaderScripts, ScriptLink{
			Src: src, IsAsync: isAsync, IsDefer: isDefer, Type: scriptType, Integrity: integrity, CrossOrigin: crossOrigin,
		})
	}
}

// WithCommonLibInclusion flags whether to include common JS libraries.
func WithCommonLibInclusion(includeHTMX, includeHTMXPreload, includeAlpine bool) Option {
	return func(vm *HeadViewModel) {
		vm.IncludeHTMX = includeHTMX
		vm.IncludeHTMXPreload = includeHTMXPreload
		vm.IncludeAlpineJS = includeAlpine
	}
}

// WithCommonLibPaths overrides default paths for common JS libraries.
func WithCommonLibPaths(htmxPath, htmxPreloadPath, alpinePath string) Option {
	return func(vm *HeadViewModel) {
		if htmxPath != "" {
			vm.HTMXPath = htmxPath
		}
		if htmxPreloadPath != "" {
			vm.HTMXPreloadPath = htmxPreloadPath
		}
		if alpinePath != "" {
			vm.AlpineJSPath = alpinePath
		}
	}
}

// WithCustomMetaTag adds a custom meta tag.
func WithCustomMetaTag(key, value string) Option {
	return func(vm *HeadViewModel) {
		vm.CustomMetaTags[key] = value // Map was initialized in NewHeadViewModel
	}
}
