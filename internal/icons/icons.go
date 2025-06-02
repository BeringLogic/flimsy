package icons

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)


func DownloadIcon(filename string) error {
  if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
    return fmt.Errorf("Invalid filename: %s", filename)
  }

  extension := strings.Replace(path.Ext(filename), ".", "", 1)
  source := "https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/" + extension + "/" + filename
  destination := "/data/icons/" + filename

  _, err := os.Stat("/data/icons"); if err != nil {
    if os.IsNotExist(err) {
      if err := os.Mkdir("/data/icons", 0755); err != nil {
        return fmt.Errorf("Could not create /data/icons directory: %w", err)
      }
    } else {
      return fmt.Errorf("os.Stat error: %w", err)
    }
  }

  bytes, err := findIcon(source); if err != nil {
    _, err := os.Stat(destination); if err == nil {
      return nil
    }
    return fmt.Errorf("Could not find icon: %s. Available icons: https://github.com/homarr-labs/dashboard-icons/blob/main/ICONS.md\n%w", filename, err)
  }

  f, err := os.Create(destination); if err != nil {
    return fmt.Errorf("Could not create /data/icons/%s: %w", filename, err)
  }
  defer f.Close()

  if _, err = f.Write(bytes); err != nil {
    return fmt.Errorf("Could not write to /data/icons/%s: %w", filename, err)
  }

  return nil
}

func findIcon(source string) ([]byte, error) {

  response, err := http.Get(source); if err != nil {
    return nil, err
  }
  defer response.Body.Close()

  if response.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("Failed to fetch icon. Status code: %d", response.StatusCode)
  }

  bytes, err := io.ReadAll(response.Body); if err != nil {
    return nil, err
  }

  return bytes, nil
} 
