package middleware


import (
  "time"
  "net/http"

  "github.com/BeringLogic/flimsy/internal/logger"
)


type wrappedWriter struct {
  http.ResponseWriter
  statusCode int
}


func (w *wrappedWriter) WriteHeader(statusCode int) {
  w.ResponseWriter.WriteHeader(statusCode)
  w.statusCode = statusCode
}

func logging(log *logger.FlimsyLogger, logAll bool, next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    wrapped := &wrappedWriter{
      ResponseWriter: w,
      statusCode: http.StatusOK,
    }

    next.ServeHTTP(wrapped, r)

    if logAll || wrapped.statusCode > 399 {
      log.Printf("%d | %s | %s | %s", wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
    }
  })
}

func Logging(log *logger.FlimsyLogger, logAll bool) func(http.Handler) http.Handler {
  return func(next http.Handler) http.Handler {
    return logging(log, logAll, next)
  }
}
