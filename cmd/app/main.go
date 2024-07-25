package main

import (
	"messagio_testsuite/internal/app"
)

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
