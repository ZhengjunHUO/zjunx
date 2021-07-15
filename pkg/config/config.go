package config

import (
	"log"
	"encoding/json"
	"io/ioutil"
)

// Modify the plafond according to your infrastructure
const MAX_CONN_NUM uint64 = 65535
const MAX_WORKER_NUM uint64 = 2048
const MAX_QUEUE_SIZE uint64 = 128

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
	Cfg = &Config{}
	if !checkConfig() {
		Cfg = &Config {
			ServerName: "Zjunx Server",
			ListenIP: "127.0.0.1",
			ListenPort: 8080,
			ConnLimit: 128,
			WorkerProcesses: 1,
			BacklogSize: 1,
			ScheduleAlgo: "RoundRobin",
		}
	}
}

func checkConfig() bool {
	// Read user defined configuration file
	content, err := ioutil.ReadFile("../../config/zjunx.cfg")
	// if user defined conf not found use the default config
	if err != nil {
		log.Println("[WARN] Unable to read the config file: ", err)
		return false
	}

	// syntax errors should be corrected in config file before ZJunx server runs
	if err = json.Unmarshal(content, Cfg); err != nil {
		log.Fatalln("[FATAL] Error occurred when parsing config file: ", err)
	}

	// regularize user defined max values
	switch {
		case Cfg.ConnLimit < 1:
			Cfg.ConnLimit = 1
		case Cfg.ConnLimit > MAX_CONN_NUM:
			Cfg.ConnLimit = MAX_CONN_NUM
	}

	switch {
		case Cfg.WorkerProcesses < 1:
			Cfg.WorkerProcesses = 1
		case Cfg.WorkerProcesses > MAX_WORKER_NUM:
			Cfg.WorkerProcesses = MAX_WORKER_NUM
	}

	switch {
		case Cfg.BacklogSize < 1:
			Cfg.BacklogSize = 1
		case Cfg.BacklogSize > MAX_QUEUE_SIZE:
			Cfg.BacklogSize = MAX_QUEUE_SIZE
	}

	log.Println("[Info] Config file loaded.")
	return true
}
