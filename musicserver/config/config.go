package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	VidFolder       string
	VidExe          string
	VidArgs         []string
	VidTimout       string
	AdminPass       string
	AdminPassHashed string
	ServerDomain    string
	TemplateDir     string
	YTDLExe         string
	FFMPEGExe       string
	Port            string
	Buckets         int
}

func ReadConfig() (Config, error) {
	out := Config{}

	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		return out, err
	}

	err = json.Unmarshal(bytes, &out)
	if err != nil {
		return out, err
	}

	return out, nil
}
