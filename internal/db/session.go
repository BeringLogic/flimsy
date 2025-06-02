package db

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)


type Session struct {
  Token string
  ExpiresAt time.Time
}


func (token *Session) IsExpired() bool {
  return time.Now().UTC().After(token.ExpiresAt)
}


func (db *FlimsyDB) DeleteExpiredSessions() error {
  _, err := db.sqlDb.Exec("DELETE FROM session WHERE expires_at < ?", time.Now().Format(time.DateTime))
  return err
}

func (db *FlimsyDB) LoadSessions() ([]*Session, error) {
  Sessions := make([]*Session, 0);

  rows, err := db.sqlDb.Query("SELECT * FROM session"); if err != nil {
    return nil, err
  }

  expiresAtString := "" 
  for rows.Next() {
    var session Session
    if err = rows.Scan(&session.Token, &expiresAtString); err != nil {
      return nil, err
    }
    session.ExpiresAt, err = time.Parse(time.DateTime, expiresAtString); if err != nil {
      return nil, err
    }
    Sessions = append(Sessions, &session)
  }

  return Sessions, nil
}

func generateToken() (string, error) {
  bytes := make([]byte, 32)
  if _, err := rand.Read(bytes); err != nil {
    return "", err
  }

  token := base64.URLEncoding.EncodeToString(bytes)
  return token, nil
}

func (db *FlimsyDB) CreateNewSession() (*Session, error) {
  var err error
  session := new(Session)
  session.Token, err = generateToken(); if err != nil {
    return nil, err
  }
  session.ExpiresAt = time.Now().UTC().Add(time.Hour * 3)

  _, err = db.sqlDb.Exec("INSERT INTO session (token, expires_at) VALUES (?, ?)", session.Token, session.ExpiresAt.Format(time.DateTime)); if err != nil {
    return nil, err
  }

  return session, nil
}

func (db *FlimsyDB) DeleteSession(token string) error {
  _, err := db.sqlDb.Exec("DELETE FROM session WHERE token = ?", token)
  return err
}
