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
func getDbPath() string {
	dataHome := getDataHomePath()
	dataAppFolder := filepath.Join(dataHome, "gotop")

	if err := os.MkdirAll(dataAppFolder, 0700); err != nil {
		log.Fatalf("Failed to create data dir at %s: %s", dataAppFolder, err)
	}

	return filepath.Join(dataAppFolder, "db.sqlite3")
}

// Connect the database to create a db instance
func getDb(resetDb bool) db {
	dbPath := getDbPath()
	if resetDb {
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
	UNIQUE(pid, start)
)`

// Ensure database is initialised with the right tables
func initDb(db *sql.DB) {
	if _, err := db.Exec(createProcessTableQuery); err != nil {
		log.Fatalf("[sql error] Failed to create processes table: %s", err)
	}
}

const insertProcessQuery string = `
INSERT INTO processes(pid, start, end, cmdline) VALUES(?, ?, ?, ?)
ON CONFLICT(pid, start) DO UPDATE SET
	end = excluded.end,
	cmdline = excluded.cmdline;`

// Store the current process
func storeProcesses(processes *processList, db *db) {
	for {
		time.Sleep(time.Second)

		processes_view := copyProcesses(processes)

		for _, process := range processes_view {
			db.mutex.Lock()

			_, err := db.db.Exec(insertProcessQuery, process.pid, process.start, process.end, process.cmdline)

			db.mutex.Unlock()

			if err != nil {
				log.Fatalf("Failed to update process %s: %s", process.cmdline, err)
			}
		}
	}
}
