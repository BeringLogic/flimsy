# TODO
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
* [x] Support webp icons
* [x] SQLite
* [x] Customizable icon in h1
* [x] Customizable title in h1
* [x] Customizable background image
* [x] Customizable number of columns (customizable per list)
* [x] Customizable colors 
* [x] File upload for backgrounds
* [x] Autodetect the colors from the background image
* [x] Catppuccin colors (mocha and latte)
* [x] Output of lm-sensors in a list so users can figure out their CPU sensor's name
* [x] Server status (CPU temp)
* [x] Server status (RAM)
* [x] Server status (SWAP)
* [x] Server status (Storage)
* [x] Server status (public IP)
* [x] Make RAM, SWAP, public ip and storage optional (/ is there by default)
* [x] Display online status
* [ ] UI to customize login variables
* [ ] UI to customize weather widget variables
* [ ] Update README.md with screenshots, features, infos on how to install and how to use the app
* [ ] Docker integration
* [ ] favicons for items (https://dashy.to/docs/icons/)
* [x] HTMX ?
* [x] Update system info and weather every X minutes
* [ ] /signup to create user
* [ ] Keyboard shortcuts?
* [ ] RSS feeds?
* [ ] Pihole widget?

# Switch to golang
* [x] golang base image
* [x] hello world 
* [x] golang baseimage example
* [x] watch source or bind src/ ?
* [x] compile automatically?
* [ ] CI pipeline?
* [x] net/http server
* [x] serve static files
* [x] /
* [x] repository pattern
* [x] login
* [x] sqlite
* [x] debugger
* [ ] Dockerfile app prod (from scratch)
* [ ] tests?
* [ ] ADDRESS, INTERFACE? or DOMAIN?
* [ ] PORT
* [x] Logger
* [ ] Let's Encrypt?
* [ ] Use the 404 template

# Docker stuff I need to learn
* [ ] Docker secrets
* [ ] How to deal with /data. I want it to work with a managed volume, but I want it to work with a bound volume too.
* [ ] How to run rootless (https://docs.docker.com/engine/security/rootless/)
* [ ] Find the proper way to deal with Dev and Production.

# Bugs
- http: superfluous response.WriteHeader call from github.com/BeringLogic/flimsy/internal/middleware.(*wrappedWriter).WriteHeader (logging.go:19)
  I get this if I do a w.Write() before w.WriteHeader() but not if I do http.Error()? Looks like flimsyServer.executeTemplate is doing it too
  That call isn't superfluous, it's to set the error code
- Error when clicking a list or an item and also when drag&dropping: "TypeError: this is undefined" in moz-extension://c1bc3b73-63a9-49a2-8dda-77940ea6c9de/content/bootstrap-legacy-autofill-overlay.js
  Appears to be a Bitwarden and Sortable.js issue
- mounted folders in /mnt are listed alphabetically, not in the order in the compose file
- if icon download fails but there is an icon of that name, just use it
- /onlinestatus should timeout after 3 seconds
- /onlinestatus shows offline for glances, pi-hole, portainer, dockge, all the syncthings, and router, even though they are online

