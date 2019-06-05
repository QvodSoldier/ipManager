package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// io.WriteString(w, config.Cfg.Url)

type Config struct {
	Nsxmanager []map[string]string
	Username   string
	Password   string
}

var Cfg = &Config{}

func Init() {
	log.SetPrefix("[ipmanager]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	data, err := ioutil.ReadFile("/alauda/config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &Cfg)
	if err != nil {
		log.Fatal(err)
	}
}
