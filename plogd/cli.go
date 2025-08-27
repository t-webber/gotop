package main

import (
	"fmt"
	"log"
	"os"
)

import "plog/lib"

// Read given argvs to build args
func parseArgs() lib.Args {
	args := lib.Args{ResetDb: false, Display: false, DbPath: ""}
	waitingForPath := false

	for idx, elt := range os.Args {
		if idx == 0 {
			continue
		}

		if waitingForPath {
			args.DbPath = elt
			waitingForPath = false
			continue
		}

		switch elt {

		case "--display":
			args.Display = true
		case "--resetdb":
			args.ResetDb = true
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
