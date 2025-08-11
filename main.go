package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"os"
)

type Process struct {
	mutex   sync.Mutex
	cmdline string
	start   time.Time
	end     time.Time
}

type ProcessList struct {
	mutex sync.Mutex
	list  map[int]*Process
}

func updateProcessesFromList(processes *ProcessList, files []os.DirEntry) {

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		pid := file.Name()
		pid_int, err := strconv.Atoi(pid)
		if err != nil {
			continue
		}

		cmdlinePath := filepath.Join("/proc", pid, "cmdline")
		data, err := os.ReadFile(cmdlinePath)
		if err != nil || len(data) == 0 {
			continue
		}
		cmdline := strings.ReplaceAll(string(data), "\x00", " ")

		processes.mutex.Lock()
		value, ok := processes.list[pid_int]
		if ok {
			processes.mutex.Unlock()
			value.mutex.Lock()
			value.end = time.Now()
			value.mutex.Unlock()
		} else {
			processes.list[pid_int] = &Process{end: time.Now(), start: time.Now(), cmdline: cmdline}
			processes.mutex.Unlock()
		}
	}
}

func updateProcesses(processes *ProcessList) {
	for {
		files, err := os.ReadDir("/proc")
		if err != nil {
			log.Fatal(err)
		}
		updateProcessesFromList(processes, files)

	}
}

func displayProcesses(processes *ProcessList) {
	for {
		processes.mutex.Lock()
		for pid, process := range processes.list {
			process.mutex.Lock()
			if strings.Contains(process.cmdline, "alacritty") {
				fmt.Printf("!%d: %s (%s->%s)!\n", pid, process.cmdline, process.start, process.end)
			}
			process.mutex.Unlock()
		}
		processes.mutex.Unlock()
	}
}

func main() {
	processes := ProcessList{list: make(map[int]*Process)}

	go updateProcesses(&processes)
	go displayProcesses(&processes)

	select {}
}
