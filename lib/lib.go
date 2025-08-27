package lib

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
)

// Represents the arguments passed by the user through argv
type Args struct {
	ResetDb bool
	Display bool
	DbPath  string
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
func getDbPath(argsDbPath string) string {
	if argsDbPath != "" {
		return argsDbPath
	}
	dataHome := getDataHomePath()
	dataAppFolder := filepath.Join(dataHome, "gotop")

	if err := os.MkdirAll(dataAppFolder, 0700); err != nil {
		log.Fatalf("Failed to create data dir at %s: %s", dataAppFolder, err)
	}

	return filepath.Join(dataAppFolder, "db.sqlite3")
}

// Connect the database to create a db instance
func GetDb(args Args) *sql.DB {
	dbPath := getDbPath(args.DbPath)
	if args.ResetDb {
		if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
			log.Fatalf("Failed to remove %s: %s", dbPath, err)
		}
	}
	log.Printf("Saving data to %s.\n", dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Access to %s denied: %s", dbPath, err)
	}
	return db
}
