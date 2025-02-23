package systemInfo


import (
  "strings"
  "os/exec"
)


type Chip struct {
  Name string
  Sensors []*Sensor
}

type Sensor struct {
  Name string
  Readings []*Reading
}

type Reading struct {
  Name string
  Value string
}

func GetSensors() ([]*Chip, error) {
  output, err := exec.Command("sensors", "-uA").Output(); if err != nil {
    return nil, err
  }

  sensors := []*Chip{}
  var chip *Chip
  var sensor *Sensor
  var reading *Reading

  for _, line := range strings.Split(string(output), "\n") {
    line = strings.TrimSpace(line)

    if line == "" {
      continue
    }

    if !strings.Contains(line, ":") {
      chip = new(Chip)
      chip.Name = line
      chip.Sensors = []*Sensor{}
      sensors = append(sensors, chip)
      continue
    }

    parts := strings.Split(line, ":")
    if parts[1] == "" {
      sensor = new(Sensor)
      sensor.Name = parts[0]
      sensor.Readings = []*Reading{}
      chip.Sensors = append(chip.Sensors, sensor)
    } else {
      reading = new(Reading)
      reading.Name = parts[0]
      reading.Value = parts[1]
      sensor.Readings = append(sensor.Readings, reading)
    }
  }
  
  return sensors, nil
}

