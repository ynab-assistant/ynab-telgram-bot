package main

import (
	"github.com/oneils/ynab-helper/bot/pkg/app"
)

const configPath = "configs/main"

func main() {
	app.Run(configPath)
}
