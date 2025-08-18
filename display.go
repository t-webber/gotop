package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

type processDisplay struct {
	pid     int
	start   int
	end     time.Time
	cmdline string
}

func displayProcesses(processes *processList) {
	for {
		time.Sleep(time.Second)

		processes_view := []processDisplay{}

		processes.mutex.Lock()
		for id, process := range processes.list {
			process.mutex.Lock()
			processes_view = append(processes_view, processDisplay{pid: id.pid, start: id.start, cmdline: process.cmdline, end: process.end})
			process.mutex.Unlock()
		}
		processes.mutex.Unlock()

		sort.Slice(processes_view, func(i, j int) bool {
			return processes_view[i].pid < processes_view[j].pid

		})

		fmt.Print("\033[3J\033[H\033[2J")
		for _, process := range processes_view {

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
}
