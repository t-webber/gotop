package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

type processDisplay struct {
	pid      int
	cmdline  string
	duration int64
}

// Copy the processList to unlock the mutex faster
func getDisplayProcesses(processes *processList, btime int64) []processDisplay {

	processes_view := []processDisplay{}

	processes.mutex.Lock()
	for id, process := range processes.list {
		process.mutex.Lock()
		duration := process.end.Unix() - id.start - btime
		processes_view = append(processes_view, processDisplay{pid: id.pid, duration: duration, cmdline: id.cmdline})
		process.mutex.Unlock()
	}
	processes.mutex.Unlock()

	return processes_view
}

// Display the list of processes
func displayProcessSync(processes []processDisplay) {
	for _, process := range processes {
		if process.cmdline == "" {
			continue
		}

		cmdpath := strings.SplitN(process.cmdline, " ", 2)[0]
		path_elts := strings.Split(cmdpath, "/")
		cmdprog := path_elts[len(path_elts)-1]

		if cmdprog == "" && process.cmdline != "" {
			log.Fatalf("%s generated empty command prog (path = %s)", process.cmdline, cmdpath)
		}

		const timeFormat = "2006/01/02 15:04:05"

		fmt.Printf("%05d %-25s %02d:%02d:%02d\n",
			process.pid,
			cmdprog,
			process.duration/3600,
			process.duration%3600/60,
			process.duration%60,
		)
	}
}

// Copies the processes, sorts them and displays them
func displayProcesses(processes *processList, btime int64) {
	for {
		time.Sleep(time.Second)

		processes_view := getDisplayProcesses(processes, btime)

		sort.Slice(processes_view, func(i, j int) bool {
			x := processes_view[i]
			y := processes_view[j]
			return (x.duration > y.duration) || (x.duration == y.duration && x.pid > y.pid)
		})

		fmt.Print("\033[3J\033[H\033[2J")
		fmt.Printf("%-5s %-25s %-8s\n", "pid", "cmdline", "duration")
		fmt.Printf(strings.Repeat("-", 40) + "\n")
		displayProcessSync(processes_view)
	}
}
