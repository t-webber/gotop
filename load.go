package main

import (
	"database/sql"
	"log"
	"time"
)

func loadProcesses(db *sql.DB, processes *processList) {
	rows, err := db.Query("SELECT pid, start, cmdline, end FROM processes")
	if err != nil {
		log.Fatalf("[sql error] Failed to load processes: %s", err)
	}
	defer rows.Close()

	var id processId
	var end time.Time

	for rows.Next() {
		if err := rows.Scan(&id.pid, &id.start, &id.cmdline, &end); err != nil {
			log.Fatalf("[sql error] Invalid row in processes load: %s", err)
		}

		processes.mutex.Lock()
		if _, ok := processes.list[id]; !ok {
			processes.list[id] = &process{end: end}

		}
		processes.mutex.Unlock()
	}
}
