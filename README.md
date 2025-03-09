# Flimsy
Flimsy home page for your homelab

## Features
- Customizable icon, title, background image and colors
- Use dashboard-icons from homarr-labs
- Catppuccin Latte and Mocha themes
- Autodetection of colors from the uploaded background image
- System-information including CPU temperature, free memory and swap, public IP and free space on filesystems
- Current Weather from OpenWeatherMap

## Installation

### Create a compose.yaml file
```yaml
services:
  app:
    image: beringlogic/flimsy-app:latest
    volumes:
      - data:/data
            
      # If you want to monitor free space on a filesystem, mount it read-only on /mnt/name in the container
      - /home/phil/Data:/mnt/Data:ro
      - /mnt/backups:/mnt/backups:ro
    environment:
      # Optional: Only used for log files...
      - TZ=${TZ:-UTC}

      # Optional: Login credentials. If none are specified, authentication is disabled.
      - FLIMSY_USERNAME=${FLIMSY_USERNAME}
      - FLIMSY_PASSWORD=${FLIMSY_PASSWORD}

      # Optional: OpenWeatherMap.org API settings.
      - FLIMSY_WEATHER_API_KEY=${FLIMSY_WEATHER_API_KEY}
      - FLIMSY_WEATHER_LOCATION=${FLIMSY_WEATHER_LOCATION}
      - FLIMSY_WEATHER_UNITS=${FLIMSY_WEATHER_UNITS}
      - FLIMSY_WEATHER_LANGUAGE=${FLIMSY_WEATHER_LANGUAGE}

      # Optional: Name of the CPU Temp sensor in lm-sensors. ex: "Tctl" for Ryzen CPUs or "Package id 0" for Xeon CPUs.
      - FLIMSY_CPU_TEMP_SENSOR=${FLIMSY_CPU_TEMP_SENSOR}
```

### Start the container
```bash
docker compose up -d
```

## Usage
- When logged out, click on items to go there
- Login to be able to edit
- Click the gear button to customize appearance
- Click the lists and items to edit them
- Drag & drop to reorder lists and items

## References
- [catppuccin themes](https://github.com/catppuccin/catppuccin/blob/main/docs/style-guide.md)
- [dashboard-icons](https://github.com/homarr-labs/dashboard-icons)
- The favicon is [homarr-labs/dashboard-icons homepage.png](https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png/homepage.png)
- The autodetection of colors from the uploaded background is using this code: [palette-extractor](https://github.com/BeringLogic/palette-extractor)

