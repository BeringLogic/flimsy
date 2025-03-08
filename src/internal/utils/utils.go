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

func GetBackgrounds() ([]string, error) {
  files, err := os.ReadDir("/data/backgrounds"); if err != nil {
    return []string{}, err
  }

  filenames := []string{}
  for _, file := range files {
    filenames = append(filenames, file.Name())
  }

  return filenames, nil
}
