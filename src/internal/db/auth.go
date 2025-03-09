package db

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)


type AuthTokenPair struct {
  Id int64
  SessionToken string
  CsrfToken string
  ExpiresAt time.Time
}


func (tokenPair *AuthTokenPair) IsExpired() bool {
  return time.Now().UTC().After(tokenPair.ExpiresAt)
}

func (db *FlimsyDB) LoadAuthTokens() ([]*AuthTokenPair, error) {
  AuthTokenPairs := make([]*AuthTokenPair, 0);

  rows, err := db.sqlDb.Query("SELECT * FROM auth_tokens"); if err != nil {
    return nil, err
  }

  expiresAtString := "" 
  for rows.Next() {
    var token AuthTokenPair
    if err = rows.Scan(&token.Id, &token.SessionToken, &token.CsrfToken, &expiresAtString); err != nil {
      return nil, err
    }
    token.ExpiresAt, err = time.Parse(time.DateTime, expiresAtString); if err != nil {
      return nil, err
    }
    AuthTokenPairs = append(AuthTokenPairs, &token)
  }

  return AuthTokenPairs, nil
}

func generateToken() (string, error) {
  bytes := make([]byte, 32)
  if _, err := rand.Read(bytes); err != nil {
    return "", err
  }

  token := base64.URLEncoding.EncodeToString(bytes)
  return token, nil
}

func (db *FlimsyDB) GenerateTokenPair() (*AuthTokenPair, error) {
  var err error
  tokens := new(AuthTokenPair)
  tokens.SessionToken, err = generateToken(); if err != nil {
    return nil, err
  }
  tokens.CsrfToken, err = generateToken(); if err != nil {
    return nil, err
  }
  tokens.ExpiresAt = time.Now().UTC().Add(time.Hour * 3)

  result, err := db.sqlDb.Exec("INSERT INTO auth_tokens (session_token, csrf_token, expires_at) VALUES (?, ?, ?)", tokens.SessionToken, tokens.CsrfToken, tokens.ExpiresAt.Format(time.DateTime)); if err != nil {
    return nil, err
  }

  tokens.Id, err = result.LastInsertId(); if err != nil {
    return nil, err
  }

  return tokens, nil
}

func (db *FlimsyDB) DeleteExpiredTokens() error {
  _, err := db.sqlDb.Exec("DELETE FROM auth_tokens WHERE expires_at < ?", time.Now().Format(time.DateTime))
  return err
}

func (db *FlimsyDB) DeleteTokenPair(id int64) error {
  _, err := db.sqlDb.Exec("DELETE FROM auth_tokens WHERE id = ?", id)
  return err
}
