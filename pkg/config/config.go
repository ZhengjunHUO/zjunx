package config

import (
	"log"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ServerName	string
	ListenIP	string
	ListenPort	uint16

	ConnLimit	uint64

	WorkerProcesses uint64
	BacklogSize	uint64
	ScheduleAlgo	string
}

var Cfg *Config

// Default configuration values
func init() {
	Cfg = &Config {
		ServerName: "Zjunx Server",
		ListenIP: "127.0.0.1",
		ListenPort: 8080,
		ConnLimit: 256,
		WorkerProcesses: 1,
		BacklogSize: 2,
		ScheduleAlgo: "rr",
	}

	Cfg.load()
}

// Read user defined configuration file
func (c *Config) load() {
	content, err := ioutil.ReadFile("../../config/zjunx.cfg")
	if err != nil {
		log.Fatalln("Unable to read the config file: ", err)
	}

	if err = json.Unmarshal(content, c); err != nil {
		log.Fatalln("Error occurred when parsing config file: ", err)
	}
}
