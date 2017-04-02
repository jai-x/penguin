# Penguin

The music server written in Golang

### Program structure

* `admin` package
Contains data structures and functions to keep track of admin IP addresses that access the system. These IP addresses are granted permission to the admin page and the admin URL endpoints.

* `config` Package
Contains functions to read the config.json file and store it as a local struct

* `help` package
Contains some commonly used convenience functions used throughout the program

* `musicserver` package
Contains main HTTP server. Defines URL endpoints for users, admin and AJAX user requests. Defines function handlers for each URL endpoint. Initialises ProcessQueue and Admin structs and stores each instance as a variable accessed by the handler functions.

* `state` package
Defines ProcessQueue struct and functions which manage the queue of videos. Stores a Downloader instance from the `youtube` package to use to download videos. Provides the videoplayer process. Keeps a cached version of the video queue in a bucket representation. Provides functions to update and fetch the bucket representation, which is used to provide user facing information.

* `templatecache` package
Simple wrapper around the `html/template` functions to cache the templates. Templates were previously parsed and executed at runtime. This package parses and stores the parsed templates in memory at program startup. Handler functions use the `Render` function with a template name parameter.

* `youtube` package
Simple wrapper around the youtube-dl command line program. Will self-update youtube-dl on struct init.

### Features? (I don't know what to call this section)

* When returning a HTML page, the handlers in the `musicserver` will use templates found in the `templates/` directory. Makes use of the `html/template` package. Full UTF-8 is supported so you can use emoji and stuff for user aliases!

* The `musicserver` package serves the contents of the `static/` directory directly.

* The `config` package requires the JSON config file. This must be passed as the first argument to the program.

* Uses the same [ZURB Foundation CSS](http://foundation.zurb.com/) found on the [UWCS Website](http://uwcs.co.uk) for a more unified look.

* Uses [jquery](https://jquery.com/), [notify.js](https://notifyjs.com/) for the frontend webpage.

* Comes with a self-updating youtube-dl ELF binary, feel free to not use it if you don't want to, just edit the config.

### Compiling and running

* Install golang
* `cd` into top level directory
* Run `go build`
* Copy `sample.config.json` to a new file and edit
* Run program: `./penguin <path/to/new/config.json>`

### Design goals

* Use only Golang standard library, compile and run anywhere
* Be fully usable on desktop, even with JavaScript disabled

### Out of scope

* Interface designed for mobile device use

### Planned features

* R9K Mode (Video can only be played once, all videos must be unique)

### License

* MIT License
