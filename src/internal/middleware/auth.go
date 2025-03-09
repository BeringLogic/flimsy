package middleware

import (
	"context"
	"net/http"

	"github.com/BeringLogic/flimsy/internal/storage"
)


var IsAuthenticatedContextKey string = "middleware.IsAuthenticated"


func isAuthenticated(s *storage.FlimsyStorage, next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    isAuthenticated := false;

    if session_cookie, err := r.Cookie("session_token"); err == nil {
      if s.CheckSessionToken(session_cookie.Value) {
        isAuthenticated = true
      }
    }

    ctx := context.WithValue(r.Context(), IsAuthenticatedContextKey, isAuthenticated)
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}

func IsAuthenticated(flimsyStorage *storage.FlimsyStorage) func (http.Handler) http.Handler {
  return func(next http.Handler) http.Handler {
    return isAuthenticated(flimsyStorage, next)
  }
}

func MustBeAuthenticated(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    isAuthenticated := r.Context().Value(IsAuthenticatedContextKey).(bool)
    if !isAuthenticated {
      http.Error(w, "Forbidden", http.StatusForbidden)
      return
    }
    next.ServeHTTP(w, r)
  })
}

func mustHaveValidCSRFToken(s *storage.FlimsyStorage, next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    session_cookie, err := r.Cookie("session_token"); if err != nil {
      http.Error(w, "Forbidden", http.StatusForbidden)
      return
    }

    csrfToken := r.FormValue("CSRF")
    if !s.CheckCsrfToken(session_cookie.Value, csrfToken) {
      http.Error(w, "Forbidden", http.StatusForbidden)
      return
    }

    next.ServeHTTP(w, r)
  })
}

func MustHaveValidCSRFToken(flimsyStorage *storage.FlimsyStorage) func (http.Handler) http.Handler {
  return func(next http.Handler) http.Handler {
    return mustHaveValidCSRFToken(flimsyStorage, next)
  }
}
