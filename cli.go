package main

import (
	"fmt"
	"log"
	"os"
)

// Represents the arguments passed by the user through argv
type args struct {
	resetDb bool
	display bool
	dbPath  string
}

// Read given argvs to build args
func parseArgs() args {
	args := args{resetDb: false, display: false, dbPath: ""}
	waitingForPath := false

	for idx, elt := range os.Args {
		if idx == 0 {
			continue
		}

		if waitingForPath {
			args.dbPath = elt
			waitingForPath = false
			continue
		}

		switch elt {

		case "--display":
			args.display = true
		case "--resetdb":
			args.resetDb = true
		case "--path":
			waitingForPath = true
		case "--help":
			fmt.Println("Usage: gotop [--path <path>][--resetdb][--display][--help]")
			os.Exit(0)

		default:
			log.Fatalf("Invalid command line argument: %s", elt)
		}

	}

	if waitingForPath {
		log.Fatalln("Missing argument for --path")

	}
	return args
}
