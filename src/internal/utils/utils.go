package utils


import (
	"os"
)


type SelectOption struct {
  Value string
  Label string
  Selected bool
}


func GetEnv(key, def string) string {
  if value, exists := os.LookupEnv(key); !exists {
    return def
  } else {
    return value
  }
}

func GetBackgrounds() []string {
  files, err := os.ReadDir("/data/backgrounds")
  if err != nil {
    return []string{}
  }

  filenames := []string{}
  for _, file := range files {
    filenames = append(filenames, file.Name())
  }

  return filenames
}
