package utils


import (
  "os"
)


func GetEnv(key, def string) string {
  if value, exists := os.LookupEnv(key); !exists {
    return def
  } else {
    return value
  }
}
