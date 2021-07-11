package config

import (
	"log"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ServerName	string
	// IP and port the ZJunx server listens on
	ListenIP	string
	ListenPort	uint16

	// Maximum active connections allowed 
	ConnLimit	uint64

	// Maximum workers who handle the request  
	WorkerProcesses uint64
	// The lenth of request queue for each worker
	BacklogSize	uint64
	// Schduling Algorithm applied for workers
	// legitime value: RoundRobin, Random, LeastConn
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
		ScheduleAlgo: "RoundRobin",
	}

	Cfg.load()
}

// Read user defined configuration file
func (c *Config) load() {
	content, err := ioutil.ReadFile("../../config/zjunx.cfg")
	// if user defined conf not found use the default config
	if err != nil {
		log.Println("[WARN] Unable to read the config file: ", err)
		return
	}

	// errors should be corrected in config file before ZJunx server runs
	if err = json.Unmarshal(content, c); err != nil {
		log.Fatalln("[FATAL] Error occurred when parsing config file: ", err)
	}
}
