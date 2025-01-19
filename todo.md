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
* [x] Automatically pick the colors from the background image
* [x] Catppuccin colors (mocha and latte)
* [x] Output of lm-sensors in a list so users can figure out their CPU sensor's name
* [x] Server status (CPU temp)
* [x] Server status (RAM)
* [x] Server status (SWAP)
* [x] Server status (Storage)
* [ ] Server status (Network)
* [ ] Server status (public IP)
* [ ] Customizable login env variables overrides
* [ ] Customizable weather widget env variables overrides
* [ ] Update README.md with features, infos on how to install and how to use the app
* [ ] Docker integration

# Stuff I need to learn
* [ ] Docker secrets
* [ ] How to run rootless (watch doesn't work: user doesn't exist)
* [ ] Find the proper way to deal with Dev and Production
* [ ] Color in errors in Docker logs?
* [ ] xdebug

# Bugs
* [ ] Make the whole item clickable (even when logged in)
* [ ] Login feedback to display AFTER the list has loaded
* [ ] Mount points cannot be monitored (except /) as they are not mounted in the container
* [ ] Error if icon doesn't exist in the dashboard icons project should be handled gracefully
* [ ] Handle errors while creating flimsy.db
* [ ] Duplication in logs
