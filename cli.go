package main

import (
	"log"
	"os"
)

// Represents the arguments passed by the user through argv
type args struct {
	resetDb bool
	display bool
}

// Read given argvs to build args
func parseArgs() args {
	args := args{resetDb: false, display: false}

	for idx, elt := range os.Args {
		if elt == "display" {
			args.display = true
		} else if elt == "resetDb" {
			args.resetDb = true
		} else if idx != 0 {
			log.Fatalf("Invalid command line argument: %s", elt)
		}
	}

	return args
}
