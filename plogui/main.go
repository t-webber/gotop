package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/t-webber/plog/lib"
)

type weightedData struct {
	weight int64
	data   string
}

func main() {
	args := lib.ParseArgs()
	db := lib.GetDb(args)
	data := getNvimUsage(db)
	displayData("Nvim", data)

}
