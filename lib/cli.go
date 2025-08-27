package lib

import (
	"fmt"
	"log"
	"os"
)

// Represents the arguments passed by the user through argv
type Args struct {
	ResetDb bool
	Display bool
	DbPath  string
}

// Read given argvs to build args
func ParseArgs() Args {
	args := Args{ResetDb: false, Display: false, DbPath: ""}
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
			fmt.Println("Usage: plogd [--path <path>][--resetdb][--display][--help]")
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
