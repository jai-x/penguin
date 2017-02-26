package main


import (
	"flag"
	"./src/musicserver"
)

func main() {
	debug := flag.Bool("debug", false, "Controls if the music server is put into debug mode")
	flag.Parse()

	musicserver.Init(*debug)
	musicserver.Run()
}