package main

import "time"

// Contains all the data that represents a prices
type processDisplay struct {
	pid     int
	start   int
	end     time.Time
	cmdline string
}

// Copy the processList to unlock the mutex faster
func copyProcesses(processes *processList) []processDisplay {

	processes_view := []processDisplay{}

	processes.mutex.Lock()
	for id, process := range processes.list {
		process.mutex.Lock()
		processes_view = append(processes_view, processDisplay{pid: id.pid, start: id.start, cmdline: id.cmdline, end: process.end})
		process.mutex.Unlock()
	}
	processes.mutex.Unlock()

	return processes_view
}
