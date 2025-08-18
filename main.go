package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	args := parseArgs()

	db := getDb(args.resetDb)
	initDb(db.db)

	processes := processList{list: make(map[processId]*process)}

	go updateProcesses(&processes)
	if args.display {
		go displayProcesses(&processes)
	} else {
		log.Println("Running...")
	}

	loadProcesses(db.db, &processes)
	go storeProcesses(&processes, &db)

	select {}
}
