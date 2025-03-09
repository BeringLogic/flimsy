package db


import (
  // "database/sql"
  "time"
)


type CsrfToken struct {
  Session string
  Token string
  ExpiresAt time.Time
}


func (token *CsrfToken) IsExpired() bool {
  return time.Now().UTC().After(token.ExpiresAt)
}

