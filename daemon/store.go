package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Wrapper for sql.DB for concurrency management
type db struct {
	db    *sql.DB
	mutex sync.Mutex
}

// Contains all the data that represents a prices
type processDb struct {
	pid     int
	start   int64
	end     int64
	cmdline string
}

// Copy the processList to unlock the mutex faster
func getDbProcesses(processes *processList) []processDb {

	processes_view := []processDb{}

	processes.mutex.Lock()
	for id, process := range processes.list {
		process.mutex.Lock()
		processes_view = append(processes_view, processDb{pid: id.pid, start: id.start, cmdline: id.cmdline, end: process.end})
		process.mutex.Unlock()
	}
	processes.mutex.Unlock()

	return processes_view
}

// Returns the value of XDG_DATA_HOME
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

// Returns the path at which the sqlite file is stored
func getDbPath(argsDbPath string) string {
	if argsDbPath != "" {
		return argsDbPath
	}
	dataHome := getDataHomePath()
	dataAppFolder := filepath.Join(dataHome, "gotop")

	if err := os.MkdirAll(dataAppFolder, 0700); err != nil {
		log.Fatalf("Failed to create data dir at %s: %s", dataAppFolder, err)
	}

	return filepath.Join(dataAppFolder, "db.sqlite3")
}

// Connect the database to create a db instance
func getDb(args args) db {
	dbPath := getDbPath(args.dbPath)
	if args.resetDb {
		if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
			log.Fatalf("Failed to remove %s: %s", dbPath, err)
		}
	}
	log.Printf("Saving data to %s.\n", dbPath)
	handle, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Access to %s denied: %s", dbPath, err)
	}
	return db{db: handle}
}

const createProcessTableQuery string = `
CREATE TABLE IF NOT EXISTS processes (
	id 	INTEGER  PRIMARY KEY AUTOINCREMENT,
	pid     INTEGER  NOT NULL,
	start   INTEGER  NOT NULL,
	end     DATETIME NOT NULL,
	cmdline TEXT 	 NOT NULL,
	UNIQUE(pid, start, cmdline)
);

CREATE TABLE IF NOT EXISTS battery (
	time     DATETIME NOT NULL UNIQUE,
	level    INTEGER  NOT NULL,
	charging BOOLEAN  NOT NULL
);`

// Ensure database is initialised with the right tables
func initDb(db *sql.DB) {
	if _, err := db.Exec(createProcessTableQuery); err != nil {
		log.Fatalf("[sql error] Failed to create processes table: %s", err)
	}
}

const insertProcessQuery string = `
INSERT INTO processes(pid, start, cmdline, end) VALUES(?, ?, ?, ?)
ON CONFLICT(pid, start, cmdline) DO UPDATE SET
	end = excluded.end;`

// Store the current process
func storeProcesses(processes *processList, db *db) {
	for {
		time.Sleep(time.Second)

		processes_view := getDbProcesses(processes)

		db.mutex.Lock()
		tx, err := db.db.Begin()
		if err != nil {
			log.Fatalf("[sql error] Failed to begin transaction: %s", err)
		}

		for _, process := range processes_view {
			_, err := tx.Exec(insertProcessQuery, process.pid, process.start, process.cmdline, process.end)

			if err != nil {
				log.Fatalf("Failed to update process %s: %s", process.cmdline, err)
			}
		}

		err = tx.Commit()
		db.mutex.Unlock()

		if err != nil {
			log.Fatalf("[sql error] Failed to commit transaction: %s", err)
		}
	}
}
