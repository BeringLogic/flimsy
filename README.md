# Flimsy
Flimsy home page for your homelab

## Features
- Customizable icon, title, background image and colors
- Use dashboard-icons from homarr-labs
- Catppuccin Latte and Mocha themes
- Autodetection of colors from the uploaded background image
- System-information including CPU temperature, free memory and swap, public IP and free space on filesystems
- Current Weather from OpenWeatherMap
- Online status of each items
- Login to edit lists and items
- HTMX, HyperScript and Sortable.js

## Installation

### Create a compose.yaml file
```yaml
services:
  app:
    container_name: flimsy
    image: beringlogic/flimsy:latest
    volumes:
      - data:/data
            
      # If you want to monitor free space on a filesystem, mount it read-only on /mnt/name in the container
      - /home/phil/Data:/mnt/Data:ro
      - /mnt/backups:/mnt/backups:ro
    environment:
      # Optional: Only used for log files...
      - TZ=America/New_York

      # Optional: Login credentials. If none are specified, authentication is disabled.
      - FLIMSY_USERNAME=admin
      - FLIMSY_PASSWORD=admin

      # Optional: OpenWeatherMap.org API settings.
      - FLIMSY_WEATHER_API_KEY=abc123
      - FLIMSY_WEATHER_LOCATION=New York
      - FLIMSY_WEATHER_UNITS=imperial
      - FLIMSY_WEATHER_LANGUAGE=en

      # Optional: Name of the CPU Temp sensor in lm-sensors. ex: "Tctl" for Ryzen CPUs or "Package id 0" for Xeon CPUs.
      - FLIMSY_CPU_TEMP_SENSOR=Package id 0
volumes:
  data:
```

### Start the container
```bash
docker compose up -d
```

### Check the logs
```bash
docker compose logs -f
```

## Usage
- When logged out, click on items to navigate to the destination URL
- Login to be able to edit
- Click the gear button to customize appearance
- Click on Add List and Add Item
- Click on the lists and items to edit them
- Drag & drop to reorder lists and items

## References
- [catppuccin themes](https://github.com/catppuccin/catppuccin/blob/main/docs/style-guide.md)
- [dashboard-icons](https://github.com/homarr-labs/dashboard-icons)
- The favicon is [homarr-labs/dashboard-icons homepage.png](https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png/homepage.png)
- The autodetection of colors from the uploaded background is using this code: [palette-extractor](https://github.com/BeringLogic/palette-extractor)

