package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	SESSION_COOKIE_NAME                    = "session_id"
	UserKey             contextUserKeyType = "gotth_user_key"
)

type contextUserKeyType string

type SessionStore interface {
	// ExchangeSessionIDForUser returns the user object that corresponds to the sessionID
	ExchangeSessionIDForUser(ctx context.Context, sessionID string) (any, error)
	// InvalidateSession invalidates the session ID when no longer valid i.e. logout
	InvalidateSession(ctx context.Context, user any, sessionID string) error
}

// SessionCheck returns a new middleware (http.Handler) that checks whether the request has a
// session cookie and exchanges the session_id for the corresponding user.
// Use [GetUser] to retrieve the corresponding user.
// OnError is called when:
//   - [isSessionIDRequired] and the request is missing the cookie or has an invalid cookie [onFail]
//   - [isSessionIDRequired] and [ExchangeSessionIDForUser] fails
func SessionCheck(ss SessionStore, isSessionIDRequired bool, onError func(http.ResponseWriter, *http.Request, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie(SESSION_COOKIE_NAME)
			if err != nil {
				if !isSessionIDRequired {
					next.ServeHTTP(w, r)
					return
				}
				onError(w, r, err)
				return
			}

			err = sessionCookie.Valid()
			if err != nil {
				if !isSessionIDRequired {
					next.ServeHTTP(w, r)
					return
				}
				onError(w, r, err)
				return
			}

			user, err := ss.ExchangeSessionIDForUser(r.Context(), sessionCookie.Value)
			if err != nil {
				if !isSessionIDRequired {
					next.ServeHTTP(w, r)
					return
				}
				onError(w, r, err)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserKey, user)))
		})
	}
}

// InvalidateSession invalidates the sessionID of the current request and calls the next handler in
// the chain.
// It call onError when:
// - A session cookie cannot be found
// - The user object in the request context corresponding to the sessionID is null
// - SessionStore fails to invalidate the sessionID
func InvalidateSession(ss SessionStore, onError func(http.ResponseWriter, *http.Request, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie(SESSION_COOKIE_NAME)
			if err != nil {
				onError(w, r, err)
				return
			}

			user := GetUser(r.Context())
			if user == nil {
				onError(w, r, errors.New("user object is nil for an authorized request"))
				return
			}

			err = ss.InvalidateSession(r.Context(), user, sessionCookie.Value)
			if err != nil {
				onError(w, r, fmt.Errorf("failed to invalidate session. err %w", err))
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:       SESSION_COOKIE_NAME,
				Value:      sessionCookie.Value,
				Path:       "/",
				Expires:    time.Now().Add(-2 * time.Hour),
				RawExpires: "",
				HttpOnly:   true,
				SameSite:   http.SameSiteLaxMode,
			})

			next.ServeHTTP(w, r)
		})
	}
}

func GetUser(ctx context.Context) any {
	return ctx.Value(UserKey)
}
