package config

import (
	"os"
	"log"
	"io/ioutil"
	"encoding/json"
)

var (
	Config Configuration
)

type Configuration struct {
	AdminPass string

	VideoTimeout int
	MaxBuckets int

	YTDLFolder string
	YTDLBin string
}

func Init(configPath string) {
	confFile, fileErr := ioutil.ReadFile(configPath)
	if fileErr != nil {
		log.Println("File read error:", fileErr)
		os.Exit(1)
	}

	parseErr := json.Unmarshal(confFile, &Config)
	if parseErr != nil {
		log.Println("JSON parse error:", parseErr)
		log.Println("Did you pass the config file as the first argument?")
		os.Exit(1)
	}

	log.Println("Config file read!")
	log.Println(Config)
	os.Exit(0)
}

func End() {
	// config is now emptied
	Config = Configuration{}
}