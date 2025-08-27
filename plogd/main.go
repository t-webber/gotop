package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/t-webber/plog/lib"
)

const CLK_TLK = 100

func getBootTime() int64 {
	file, err := os.Open("/proc/stat")
	if err != nil {
		log.Fatalf("Failed to open /proc/stat: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "btime ") {
			continue
		}
		btime := line[6:]

		const base = 10
		const size = 64
		btime_secs, err := strconv.ParseInt(btime, base, size)
		if err != nil {
			log.Fatalf("/proc/stat has invalid content: %s should contain an integer: %s", line, err)
		}
		return btime_secs
	}
	log.Fatalln("/proc/stat has invalid content: btime line wasn't found")
	panic("")
}

func main() {
	args := lib.ParseArgs()

	db := db{db: lib.GetDb(args)}
	initDb(db.db)

	processes := processList{list: make(map[processId]*process)}

	go updateProcesses(&processes)
	if args.Display {
		btime := getBootTime()
		go displayProcesses(&processes, btime)
	} else {
		log.Println("Running...")
	}

	loadProcesses(db.db, &processes)
	go storeProcesses(&processes, &db)
	go battery(&db)

	select {}
}
