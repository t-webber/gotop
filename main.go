package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"os"
)

func main() {
	files, err := os.ReadDir("/proc")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		pid := file.Name()
		if _, err := strconv.Atoi(pid); err != nil {
			continue
		}
		cmdlinePath := filepath.Join("/proc", pid, "cmdline")
		data, err := os.ReadFile(cmdlinePath)
		if err != nil || len(data) == 0 {
			continue
		}
		cmdline := strings.ReplaceAll(string(data), "\x00", " ")
		fmt.Println(cmdline)

	}
}
