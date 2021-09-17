package config

import (
	"log"
	"encoding/json"
	"io/ioutil"
	"flag"
)

// Modify the plafond according to your infrastructure
const MAX_CONN_NUM uint64 = 65535
const MAX_WORKER_NUM uint64 = 2048
const MAX_QUEUE_SIZE uint64 = 128

const DefaultServerName string = "Zjunx Server"
const DefaultListenIP string = "0.0.0.0"
const DefaultListenPort uint64 = 8080
const DefaultConnLimit uint64 = 128
const DefaultWorkerProcesses uint64 = 1
const DefaultBacklogSize uint64 = 1
const DefaultScheduleAlgo string = "RoundRobin"

type Config struct {
	ServerName	string
	// IP and port the ZJunx server listens on
	ListenIP	string
	ListenPort	uint64

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

func init() {
	// Default configuration values
	Cfg = &Config {
		ServerName: DefaultServerName,
		ListenIP: DefaultListenIP,
		ListenPort: DefaultListenPort,
		ConnLimit: DefaultConnLimit,
		WorkerProcesses: DefaultWorkerProcesses,
		BacklogSize: DefaultBacklogSize,
		ScheduleAlgo: DefaultScheduleAlgo,
	}

	// Values from config file override the default
	loadConfig()

        // Parameter passed from command line with the highest priority
	parseConfig()

	// Check the eligibility of given parameter before launch the server
	checkConfig()
}

func parseConfig() {
	var serverName	string
	var listenIP	string
	var listenPort	uint64
	var connLimit	uint64
	var workerProcesses uint64
	var backlogSize	uint64
	var scheduleAlgo string

	flag.StringVar(&serverName, "n", DefaultServerName, "Server name")
	flag.StringVar(&listenIP, "l", DefaultListenIP, "IP server listens on")
	flag.Uint64Var(&listenPort, "p", DefaultListenPort, "IP port listens on")
	flag.Uint64Var(&connLimit, "c", DefaultConnLimit, "Max connections allowed")
	flag.Uint64Var(&workerProcesses, "w", DefaultWorkerProcesses, "Number of worker")
	flag.Uint64Var(&backlogSize, "s", DefaultBacklogSize, "Size of queue per worker")
	flag.StringVar(&scheduleAlgo, "a", DefaultScheduleAlgo, "Algorithm used to distribute job to worker, in [RoundRobin, Random, LeastConn]")
	flag.Parse()

	if serverName != DefaultServerName {
		Cfg.ServerName = serverName
	}

	if listenIP != DefaultListenIP {
		Cfg.ListenIP = listenIP
	}

	if listenPort != DefaultListenPort {
		Cfg.ListenPort = listenPort
	}

	if connLimit != DefaultConnLimit {
		Cfg.ConnLimit = connLimit
	}

	if workerProcesses != DefaultWorkerProcesses {
		Cfg.WorkerProcesses = workerProcesses
	}

	if backlogSize != DefaultBacklogSize {
		Cfg.BacklogSize = backlogSize
	}

	if scheduleAlgo != DefaultScheduleAlgo {
		Cfg.ScheduleAlgo = scheduleAlgo
	}
}

func loadConfig() {
	// Read user defined configuration file
	content, err := ioutil.ReadFile("/etc/zjunx/zjunx.cfg")
	// if user defined conf not found use the default config
	if err != nil {
		log.Println("[WARN] Unable to read the config file: ", err)
		return
	}

	// syntax errors should be corrected in config file before ZJunx server runs
	if err = json.Unmarshal(content, Cfg); err != nil {
		log.Fatalln("[FATAL] Error occurred when parsing config file: ", err)
	}

	log.Println("[Info] Config file loaded.")
}

func checkConfig() {
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

	algos := []string{"RoundRobin", "Random", "LeastConn"}
	dict := make(map[string]bool)
	for i := range algos {
		dict[algos[i]] = true
	}

	if _, ok := dict[Cfg.ScheduleAlgo]; !ok {
		Cfg.ScheduleAlgo = DefaultScheduleAlgo
	}
}
