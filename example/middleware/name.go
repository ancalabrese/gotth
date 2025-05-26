package middleware

import (
	"context"
	"net/http"
)

// GottherName is an example middleware that reads the name query param and sets a context value
// accordingly
func GottherName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		name := params.Get("name")
		ctx := context.WithValue(r.Context(), "name", name)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func GetGettherName(ctx context.Context) string {
	n, ok := ctx.Value("name").(string)
	if !ok || n == "" {
		return ""
	}
	return n
}
