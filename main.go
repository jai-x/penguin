package main


import (
	"os"
	"log"
	"./src/musicserver"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Pass config file as first argument to program")
		os.Exit(1)
	}
	musicserver.Init(os.Args[1])
	musicserver.Run()
}