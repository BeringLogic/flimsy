package middleware

import (
  "time"
  "strings"
  "context"
  "net/http"

  "github.com/BeringLogic/flimsy/internal/logger"
)


type wrappedWriter struct {
  http.ResponseWriter
  statusCode int
}

type Middleware func(http.Handler) http.Handler


func (w *wrappedWriter) WriteHeader(statusCode int) {
  w.ResponseWriter.WriteHeader(statusCode)
  w.statusCode = statusCode
}

func CreateStack(m ...Middleware) Middleware {
  return func(next http.Handler) http.Handler {
    for _, middleware := range m {
      next = middleware(next)
    }
    return next
  }
}


func Logging(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    wrapped := &wrappedWriter{
      ResponseWriter: w,
      statusCode: http.StatusOK,
    }

    next.ServeHTTP(wrapped, r)
    logger.Printf("| %d | %s | %s | %s", wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
  })
}

