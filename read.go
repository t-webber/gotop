package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type processId struct {
	pid   int
	start int
}

type process struct {
	mutex   sync.Mutex
	cmdline string
	end     time.Time
}

type processList struct {
	mutex sync.Mutex
	list  map[processId]*process
}

func getCmdLine(pid string) (string, error) {
	cmdlinePath := filepath.Join("/proc", pid, "cmdline")
	data, err := os.ReadFile(cmdlinePath)
	if err != nil || len(data) == 0 {
		return "", err
	}
	return strings.ReplaceAll(string(data), "\x00", " "), nil
}

func getStart(pid string) int {
	statPath := filepath.Join("/proc", pid, "stat")
	data, err := os.ReadFile(statPath)
	if err != nil {
		log.Fatalf("[PID %s]: cmdline is present, but not %s", pid, statPath)
	}

	fields := strings.Fields(string(data))
	if len(fields) < 22 {
		log.Fatalf("[PID %s]: found %d columns in %s, expected at least 22.", pid, len(fields), statPath)
	}

	start, err := strconv.Atoi(fields[21])
	if err != nil {
		log.Fatalf("[PID %s]: start time %s is not a valid number.", pid, fields[21])
	}
	return start
}

func updateProcessWithData(processes *processList, processId processId, cmdline string, end time.Time) {
	processes.mutex.Lock()

	value, ok := processes.list[processId]
	if ok {
		processes.mutex.Unlock()
		value.mutex.Lock()
		value.end = end
		value.mutex.Unlock()
	} else {
		processes.list[processId] = &process{end: end, cmdline: cmdline}
		processes.mutex.Unlock()
	}
}

func updateProcess(processes *processList, file os.DirEntry) {
	pid_str := file.Name()
	pid, err := strconv.Atoi(pid_str)
	if err != nil {
		return
	}

	cmdline, err := getCmdLine(pid_str)
	if err != nil {
		return
	}

	start := getStart(pid_str)
	processId := processId{pid: pid, start: start}
	end := time.Now()

	updateProcessWithData(processes, processId, cmdline, end)

}

func updateProcessFromList(processes *processList, files []os.DirEntry) {
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		updateProcess(processes, file)
	}
}

func updateProcesses(processes *processList) {
	for {
		files, err := os.ReadDir("/proc")
		if err != nil {
			log.Fatal(err)
		}
		updateProcessFromList(processes, files)

		time.Sleep(time.Second)
	}
}
