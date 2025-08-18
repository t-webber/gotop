package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

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

		end_time := process.end.Format("15:04:05")
		fmt.Printf("%05d %-25s %s\n", process.pid, cmdprog, end_time)
	}
}

// Copies the processes, sorts them and displays them
func displayProcesses(processes *processList) {
	for {
		time.Sleep(time.Second)

		processes_view := copyProcesses(processes)

		sort.Slice(processes_view, func(i, j int) bool {
			return processes_view[i].pid < processes_view[j].pid

		})

		fmt.Print("\033[3J\033[H\033[2J")
		displayProcessSync(processes_view)
	}
}
