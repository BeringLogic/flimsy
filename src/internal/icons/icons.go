package icons


import (
	"os"
	"fmt"
	"path"
	"strings"
	"net/http"
)


func DownloadIcon(filename string) error {
  if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
    return fmt.Errorf("Invalid filename: %s", filename)
  }

  extension := strings.Replace(path.Ext(filename), ".", "", 1)
  source := "https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/" + extension + "/" + filename
  destination := "/data/icons/" + filename

  response, err := http.Get(source); if err != nil {
    return fmt.Errorf("Could not find icon: %s. Available icons: https://github.com/homarr-labs/dashboard-icons/blob/main/ICONS.md\n%w", filename, err)
  }
  defer response.Body.Close()

  if response.StatusCode != http.StatusOK {
    return fmt.Errorf("Could not find icon: %s. Available icons: https://github.com/homarr-labs/dashboard-icons/blob/main/ICONS.md", filename)
  }

  _, err = os.Stat("/data/icons"); if err != nil {
    if os.IsNotExist(err) {
      if err = os.Mkdir("/data/icons", 0755); err != nil {
        return fmt.Errorf("Could not create /data/icons directory: %w", err)
      }
    } else {
      return fmt.Errorf("os.Stat error: %w", err)
    }
  }

  f, err := os.Create(destination); if err != nil {
    return fmt.Errorf("Could not create /data/icons/%s: %w", filename, err)
  }
  defer f.Close()

  if _, err = f.ReadFrom(response.Body); err != nil {
    return err
  }

  return nil
}
