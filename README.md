# Penguin

The music server written in Golang

TODO:

* Add support for account creds in youtube-dl to download from website that require an account to watch (e.g. niconico)
* Youtube link parsing
* R9K Mode (Video can only be played once, all videos must be unique)

### Features? (I don't know what to call this section)

* When returning a HTML page, the handlers in the `musicserver` package will use the `templatecache` package. This uses templates found in the `templates/` directory. It makes use of the `html/template` package. Full UTF-8 is supported so you can use emoji and stuff for user aliases!

* The `musicserver` package serves the contents of the `static/` directory directly.

* Uses the same [ZURB Foundation CSS](http://foundation.zurb.com/) found on the [UWCS Website](http://uwcs.co.uk) for a more unified look.

* Uses [jquery](https://jquery.com/), [notify.js](https://notifyjs.com/) for the frontend webpage.

* Comes with a self-updating youtube-dl ELF binary, feel free to not use it if you don't want to, just edit the config.

### Compiling and running

* Install golang
* `cd` into top level directory
* Run `go build`
* Copy `sample.config.json` to `config.json` and edit appropriately
* Run program: `./penguin`

### Design goals

* Use only Golang standard library, compile and run anywhere
* Be fully usable on desktop, even with JavaScript disabled

### Out of scope

* Interface designed for mobile device use


### License

* MIT License
