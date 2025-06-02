package middleware


import (
  "net/http"
)


type Middleware func(http.Handler) http.Handler


func CreateStack(m ...Middleware) Middleware {
  return func(next http.Handler) http.Handler {
    for _, middleware := range m {
      next = middleware(next)
    }
    return next
  }
}
