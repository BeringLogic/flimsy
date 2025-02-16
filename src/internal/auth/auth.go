package auth

import (
  "crypto/rand"
  "encoding/base64"
  // "golang.org/x/crypto/bcrypt"
)


var session_token string
var csrf_token string


func generateToken() (string, error) {
  bytes := make([]byte, 32)
  if _, err := rand.Read(bytes); err != nil {
    return "", err
  }

  token := base64.URLEncoding.EncodeToString(bytes)
  return token, nil
}

func GenerateTokens() (string, string, error) {
  var err error

  session_token, err = generateToken(); if err != nil {
    return "", "", err
  }

  csrf_token, err = generateToken(); if err != nil {
    return "", "", err
  }

  return session_token, csrf_token, nil
}

func CheckSessionToken(tokenToCheck string) bool {
  return tokenToCheck == session_token
}

func CheckCSRFToken(tokenToCheck string) bool {
  return tokenToCheck == csrf_token
}
