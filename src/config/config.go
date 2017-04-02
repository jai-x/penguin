package config

import (
	"os"
	"log"
	"io/ioutil"
	"encoding/json"
)

// Global config struct instance
var (
	Config Configuration
)

// Config struct
// The config.json file should mirror itself to this struct with given values
type Configuration struct {
	AdminPass string
	VideoTimeout int
	VideoPlayer string
	VideoPlayerArgs string
	MaxBuckets int
	DownloadFolder string
	YTDLBin string
	YTDLArgs string
	Port int
	TemplateFolder string
}

func Init(configPath string) {
	// Attempt to read file in
	confFile, fileErr := ioutil.ReadFile(configPath)
	if fileErr != nil {
		log.Println("File read error:", fileErr)
		os.Exit(1)
	}

	// Attempt to parse file as valid JSON
	parseErr := json.Unmarshal(confFile, &Config)
	if parseErr != nil {
		log.Println("JSON parse error:", parseErr)
		log.Println("Did you pass the config file as the first argument?")
		os.Exit(1)
	}

	log.Println("Config file read!")
}

func End() {
	// config is now emptied
	Config = Configuration{}
}
