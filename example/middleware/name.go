package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const (
	nameKey contextKey = "name_key"
)

// GottherName is an example middleware that reads the name query param and sets a context value
// accordingly
func GottherName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		name := params.Get("name")
		ctx := context.WithValue(r.Context(), nameKey, name)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func GetGottherName(ctx context.Context) string {
	n, ok := ctx.Value(nameKey).(string)
	if !ok || n == "" {
		return ""
	}
	return n
}
