package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

const min_check_interval int = 1
const max_check_interval int = 3600

const min_wait_before_restart int = 1
const max_wait_before_restart int = 3600

var generateLogs bool
var checkInterval, waitTime, maxAttempts int
var processName string

func checkErrAndExit(e error) {
	if e != nil {
		log.Fatal(e)
		os.Exit(1)
	}
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
	log.SetLevel(log.WarnLevel)
}

func main() {

	flag.IntVar(&checkInterval, "i", 5, "Check interval")
	flag.BoolVar(&generateLogs, "l", false, "Generate Logs")
	flag.StringVar(&processName, "p", "", "Process name")
	flag.IntVar(&maxAttempts, "r", 5, "Maximum retries")
	flag.IntVar(&waitTime, "t", 5, "Wait time before restart")

	flag.Parse()

	validateParams()
	displayParams()
}
