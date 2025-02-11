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

func LoadConfig() (*Config, error) {
  config := Config{}

  row := sqlDb.QueryRow("SELECT * FROM config WHERE id = 1")
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

