package systemInfo

import (
  "io"
  "fmt"
  "strings"
  "os/exec"
  "net/http"
  "path/filepath"
  "github.com/BeringLogic/flimsy/db"
)

type SystemInfo struct {
  Cpu_temp string
  Cpu_temp_ends_with_c bool
  Free_memory string
  Free_swap string
  Public_ip string
  Storage []Storage
}

type Storage struct {
  Mount_point string
  Free_space string
}

func GetSystemInfo(config *db.Config) SystemInfo {
  systemInfo := SystemInfo{}
  var err error
  var bytes []byte

  if (config.Cpu_temp_sensor != "") {
    cmd := fmt.Sprintf("sensors | grep \"%s\" | cut -d \":\" -f 2 | awk '{ print $1 }'", config.Cpu_temp_sensor)
    bytes, err = exec.Command("sh", "-c", cmd).Output(); if err == nil {
      systemInfo.Cpu_temp = strings.TrimSpace(string(bytes))
    }
    systemInfo.Cpu_temp_ends_with_c = strings.HasSuffix(systemInfo.Cpu_temp, "C");
  }
  if (config.Show_free_ram == 1) {
    bytes, err = exec.Command("sh", "-c", "free -h | grep Mem | awk '{ print $4 }'").Output(); if err == nil {
      systemInfo.Free_memory = strings.TrimSpace(string(bytes))
    }
  }
  if (config.Show_free_swap == 1) {
    bytes, err = exec.Command("sh", "-c", "free -h | grep Swap | awk '{ print $4 }'").Output(); if err == nil {
      systemInfo.Free_swap = strings.TrimSpace(string(bytes))
    }
  }
  if (config.Show_public_ip == 1) {
    resp, err := http.Get("https://api.ipify.org"); if err != nil {
      systemInfo.Public_ip = ""
    } else {
      bytes, err = io.ReadAll(resp.Body); if err == nil {
        systemInfo.Public_ip = strings.TrimSpace(string(bytes))
      }
    }
    resp.Body.Close()
  }
  if (config.Show_free_space == 1) {
    systemInfo.Storage = []Storage{}

    mountPoints := []string {
      "/",
    }
    matches, err := filepath.Glob("/mnt/*"); if err == nil {
      mountPoints = append(mountPoints, matches...)
    }

    for _, mountPoint := range mountPoints {
      storage := Storage{}

      storage.Mount_point = filepath.Base(mountPoint)

      cmd := fmt.Sprintf("df -h %s | tail -n 1 | awk '{ print $4 }'", mountPoint)  
      bytes, err = exec.Command("sh", "-c", cmd).Output(); if err == nil {
        storage.Free_space = strings.TrimSpace(string(bytes))
        systemInfo.Storage = append(systemInfo.Storage, storage)
      }
    }
  }
  
  return systemInfo
}
