package main

import (
	"log"
	"os"

	"github.com/oneils/ynab-helper/bot/pkg/app"
)

const configPath = "configs/main"

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	logger := log.New(os.Stdout, "BOT : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := app.Run(logger, configPath, build); err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
}
