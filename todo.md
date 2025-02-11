# TODO
* [x] Read /data/data.json
* [x] Render data.json
* [x] Add list GUI
* [x] Add item GUI
* [x] Edit list
* [x] Edit item
* [x] Delete list
* [x] Delete item
* [x] Drag & Drop to reorder list
* [x] Drag & Drop to reorder items
* [x] Make items clickable
* [x] Login to edit
* [x] Preserve icons aspect ratio
* [x] Deal with really long items
* [x] Weather widget
* [x] Default colors based on Catppucchin-mocha
* [x] Download icon when saving items and use local copy
* [x] Support both png and svg icons
* [x] SQLite
* [x] Customizable icon in h1
* [x] Customizable title in h1
* [x] Customizable background image
* [x] Customizable number of columns (customizable per list)
* [x] Customizable colors 
* [x] File upload for backgrounds
* [x] Automatically pick the colors from the background image
* [x] Catppuccin colors (mocha and latte)
* [x] Output of lm-sensors in a list so users can figure out their CPU sensor's name
* [x] Server status (CPU temp)
* [x] Server status (RAM)
* [x] Server status (SWAP)
* [x] Server status (Storage)
* [x] Server status (public IP)
* [x] Make RAM, SWAP, public ip and storage optional (/ is there by default)
* [x] curl items' href and display online status
* [ ] UI to customize login variables
* [ ] UI to customize weather widget variables
* [ ] Update README.md with screenshots, features, infos on how to install and how to use the app
* [ ] Docker integration
* [ ] Item ping icon or millisecs
* [ ] favicons for items (https://dashy.to/docs/icons/)
* [ ] check long title and lots of mount points
* [ ] Use ENV variables or the DB, not both (looking at you, CPU Temp Sensor)
* [ ] AJAX for everything?
* [ ] Refresh system-info every 5 min
* [ ] Refresh weather every hour

# Stuff I need to learn
* [ ] xdebug
* [ ] Docker secrets
* [ ] How to deal with /data. I want it to work with a managed volume, but I want it to work with a bound volume too. If bound, user must mkdir and chown 1000:1000 (flimsy user in the container). Awkward since 1000 might not be the user's id on the host...
* [ ] How to run rootless (https://docs.docker.com/engine/security/rootless/)
* [ ] Find the proper way to deal with Dev and Production. (php-production.ini, max_execution_time, grep max_input_time)

# Bugs
* [ ] Login feedback to display AFTER the list has loaded
