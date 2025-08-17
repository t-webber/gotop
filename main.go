package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type ProcessId struct {
	pid   int
	start int
}

type Process struct {
	mutex   sync.Mutex
	cmdline string
	end     time.Time
}

type ProcessList struct {
	mutex sync.Mutex
	list  map[ProcessId]*Process
}

func cmdLine(pid string) (string, error) {
	cmdlinePath := filepath.Join("/proc", pid, "cmdline")
	data, err := os.ReadFile(cmdlinePath)
	if err != nil || len(data) == 0 {
		return "", err
	}
	return strings.ReplaceAll(string(data), "\x00", " "), nil
}

func start(pid string) int {
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

func updateProcess(processes *ProcessList, pid_str string) {
	pid, err := strconv.Atoi(pid_str)
	if err != nil {
		return
	}

	cmdline, err := cmdLine(pid_str)
	if err != nil {
		return
	}

	start := start(pid_str)
	process_id := ProcessId{pid: pid, start: start}

	processes.mutex.Lock()

	value, ok := processes.list[process_id]
	if ok {
		processes.mutex.Unlock()
		value.mutex.Lock()
		value.end = time.Now()
		value.mutex.Unlock()
	} else {
		processes.list[process_id] = &Process{end: time.Now(), cmdline: cmdline}
		processes.mutex.Unlock()
	}
}

func updateProcessesFromList(processes *ProcessList, files []os.DirEntry) {

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		pid := file.Name()
		updateProcess(processes, pid)
	}
}

func updateProcesses(processes *ProcessList) {
	for {
		files, err := os.ReadDir("/proc")
		if err != nil {
			log.Fatal(err)
		}
		updateProcessesFromList(processes, files)

		time.Sleep(time.Second)
	}
}

type ProcessDisplay struct {
	pid     int
	start   int
	end     time.Time
	cmdline string
}

func displayProcesses(processes *ProcessList) {
	for {
		processes_view := []ProcessDisplay{}

		processes.mutex.Lock()
		for id, process := range processes.list {
			process.mutex.Lock()
			processes_view = append(processes_view, ProcessDisplay{pid: id.pid, start: id.start, cmdline: process.cmdline, end: process.end})
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

		time.Sleep(time.Second)
	}
}

func getDataHomePath() string {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome != "" {
		return dataHome
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to find home directory: %s", err)
	}
	return filepath.Join(home, ".local", "share")
}

func getDbPath() string {
	dataHome := getDataHomePath()
	dataAppFolder := filepath.Join(dataHome, "gotop")
	err := os.MkdirAll(dataAppFolder, 0744)
	if err != nil {
		log.Fatalf("Failed to create data dir at %s: %s", dataAppFolder, err)
	}
	return filepath.Join(dataAppFolder, "db.sqlite3")

}

func getDb() *sql.DB {
	dbPath := getDbPath()
	log.Printf("Saving data to %s.\n", dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Access to %s denied: %s", dbPath, err)
	}
	return db
}

func initDb(db *sql.DB) {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS processes (
	id 	INTEGER  PRIMARY KEY AUTOINCREMENT,
	pid     INTEGER  NOT NULL,
	start   INTEGER  NOT NULL,
	end     DATETIME NOT NULL,
	cmdline TEXT 	 NOT NULL
)
	`)

	if err != nil {
		log.Fatalf("[sql error] Failed to create processes table: %s", err)
	}
}

func main() {
	db := getDb()
	initDb(db)

	processes := ProcessList{list: make(map[ProcessId]*Process)}

	go updateProcesses(&processes)
	if len(os.Args) > 1 && os.Args[1] == "display" {
		go displayProcesses(&processes)
	}

	select {}
}
