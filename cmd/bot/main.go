package main

import (
	"github.com/oneils/ynab-helper/bot/pkg/app"
)

const configPath = "configs/main"

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	app.Run(configPath, build)
}
