package middleware


import (
  "context"
  "net/http"

  "github.com/BeringLogic/flimsy/internal/auth"
)


var IsAuthenticatedContextKey string = "middleware.IsAuthenticated"


func IsAuthenticated(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    isAuthenticated := false;

    if session_cookie, err := r.Cookie("session_token"); err == nil {
      if auth.CheckSessionToken(session_cookie.Value) {
        isAuthenticated = true
      }
    }

    ctx := context.WithValue(r.Context(), IsAuthenticatedContextKey, isAuthenticated)
    next.ServeHTTP(w, r.WithContext(ctx))
  })
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

func MustHaveValidCSRFToken(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    csrfToken := r.Header.Get("X-CSRF-TOKEN")
    if !auth.CheckCSRFToken(csrfToken) {
      http.Error(w, "Forbidden", http.StatusForbidden)
      return
    }
    next.ServeHTTP(w, r)
  })
}
