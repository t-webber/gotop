package main

import (
	"log"
	"os"
)

type args struct {
	resetDb bool
	display bool
}

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
