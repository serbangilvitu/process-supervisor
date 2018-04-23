package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	ps "github.com/mitchellh/go-ps"
	log "github.com/sirupsen/logrus"
)

const config_file = "config.json"
const min_check_interval int = 1
const max_check_interval int = 3600

const min_wait_before_restart int = 1
const max_wait_before_restart int = 3600

var generateLogs bool
var checkInterval, waitTime, maxAttempts int
var processName string
var processArgs string

var restartAttempts int = 0

type Config struct {
	Processes []Process `json:"processes"`
}

type Process struct {
	Command string `json:"command"`
	Args    string `json:"args"`
}

func checkErrAndExit(e error) {
	if e != nil {
		log.Fatal(e)
		os.Exit(1)
	}
}

func attemptRestart() {
	cmd := exec.Command(processName, processArgs)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.WithFields(log.Fields{"processName": processName,
			"output": stdoutStderr}).
			Error("Failed to start process")
	}
	restartAttempts++
	if restartAttempts > maxAttempts {
		checkErrAndExit(fmt.Errorf("Maximum restart attempts has been reached %d",
			maxAttempts))
	}
}

func findProcess() bool {
	found := false
	proc, _ := ps.Processes()
	for _, v := range proc {
		if v.Executable() == processName {
			found = true
			break
		}
	}
	return found
}

func validateParams() {
	if (checkInterval < min_check_interval) ||
		(checkInterval > max_check_interval) {
		checkErrAndExit(fmt.Errorf("Check interval(-i) must be between %d and %d",
			min_check_interval, max_check_interval))
	}
	if processName == "" {
		checkErrAndExit(fmt.Errorf("Process name(-p) cannot be empty"))
	}
	if (checkInterval < min_wait_before_restart) ||
		(checkInterval > max_wait_before_restart) {
		checkErrAndExit(fmt.Errorf("Wait time before restart(-t) must be between %d and %d",
			min_wait_before_restart, max_wait_before_restart))
	}
}

func displayParams() {
	log.WithFields(log.Fields{
		"checkInterval": checkInterval,
		"generateLogs":  generateLogs,
		"maxAttempts":   maxAttempts,
		"processName":   processName,
		"waitTime":      waitTime,
	}).Info("Parameters")
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	jsonFile, err := os.Open(config_file)
	if err != nil {
		checkErrAndExit(err)
	}
	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		checkErrAndExit(err)
	}

	var config Config
	json.Unmarshal(jsonData, &config)
	for k, v := range config.Processes {
		log.WithFields(log.Fields{"command": v.Command,
			"args": v.Args}).
			Infof("Listing process %v", k)
	}
}

func main() {

	flag.StringVar(&processArgs, "a", "", "Process arguments")
	flag.IntVar(&checkInterval, "i", 5, "Check interval")
	flag.BoolVar(&generateLogs, "l", false, "Generate Logs")
	flag.StringVar(&processName, "p", "", "Process name")
	flag.IntVar(&maxAttempts, "r", 3, "Maximum retries")
	flag.IntVar(&waitTime, "t", 10, "Wait time before restart")

	flag.Parse()

	validateParams()
	displayParams()

	if !generateLogs {
		log.SetLevel(log.FatalLevel)
	}

	for {
		if findProcess() {
			log.WithFields(log.Fields{"processName": processName}).
				Info("Process is running")
			restartAttempts = 0
			time.Sleep(time.Duration(checkInterval) * time.Second)
		} else {
			log.WithFields(log.Fields{"processName": processName}).
				Warn("Process is NOT running")
			attemptRestart()
			time.Sleep(time.Duration(waitTime) * time.Second)
		}
	}

}
