package db


type Config struct {
  Id int
  Icon string
  Title string
  Background_image string
  Color_background string
  Color_foreground string
  Color_items string
  Color_borders string
  Cpu_temp_sensor string
  Show_free_ram int
  Show_free_swap int
  Show_public_ip int
  Show_free_space int
}


func (flimsyDB *FlimsyDB) LoadConfig() (*Config, error) {
  config := Config{}

  row := flimsyDB.sqlDb.QueryRow("SELECT * FROM config WHERE id = 1")
  err := row.Scan(
      &config.Id,
      &config.Icon,
      &config.Title,
      &config.Background_image,
      &config.Color_background,
      &config.Color_foreground,
      &config.Color_items,
      &config.Color_borders,
      &config.Cpu_temp_sensor,
      &config.Show_free_ram,
      &config.Show_free_swap,
      &config.Show_public_ip,
      &config.Show_free_space); if err != nil {
    return nil, err
  }

  return &config, nil
}

func (flimsyDB *FlimsyDB) SaveConfig(config *Config) error {
  _, err := flimsyDB.sqlDb.Exec(
    `UPDATE config SET
      icon = ?,
      title = ?,
      background_image = ?,
      color_background = ?,
      color_foreground = ?,
      color_items = ?,
      color_borders = ?,
      cpu_temp_sensor = ?,
      show_free_ram = ?,
      show_free_swap = ?,
      show_public_ip = ?,
      show_free_space = ?
    WHERE id = 1`,
    config.Icon,
    config.Title,
    config.Background_image,
    config.Color_background,
    config.Color_foreground,
    config.Color_items,
    config.Color_borders,
    config.Cpu_temp_sensor,
    config.Show_free_ram,
    config.Show_free_swap,
    config.Show_public_ip,
    config.Show_free_space)

  return err
}
