package db

type Config struct {
  ID int
  icon string
  title string
  background_image string
  color_background string
  color_foreground string
  color_items string
  color_borders string
  cpu_temp_sensor string
  show_free_ram int
  show_free_swap int
  show_public_ip int
  show_free_space int
}

func LoadConfig() (*Config, error) {
  config := Config{}

  row := sqlDb.QueryRow("SELECT * FROM config WHERE id = 1")
  err := row.Scan(
      &config.ID,
      &config.icon,
      &config.title,
      &config.background_image,
      &config.color_background,
      &config.color_foreground,
      &config.color_items,
      &config.color_borders,
      &config.cpu_temp_sensor,
      &config.show_free_ram,
      &config.show_free_swap,
      &config.show_public_ip,
      &config.show_free_space); if err != nil {
    return nil, err
  }

  return &config, nil
}

